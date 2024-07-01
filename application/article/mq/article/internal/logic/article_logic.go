package logic

import (
	"beyond/application/article/model"
	"beyond/application/article/mq/article/internal/svc"
	"beyond/application/article/mq/article/internal/types"
	"beyond/application/user/rpc/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
	"time"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "mq article", Value: "article_logic"}),
	}
}

func (l *ArticleLogic) Consume(_, val string) error {
	l.Logger.Infof("Consume msg val: %s", val)
	var msg *types.CanalArticleMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		l.Logger.Errorf("Consume val: %s error: %v", val, err)
		return err
	}
	return l.articleOperate(msg)
}

// 消费文章变更数据核心逻辑(异步处理,对缓存进行补偿)
func (l *ArticleLogic) articleOperate(msg *types.CanalArticleMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}
	var esData []*types.ArticleEsMsg
	for _, d := range msg.Data {
		status, _ := strconv.Atoi(d.Status)
		likeNum, _ := strconv.ParseInt(d.LikeNum, 10, 64)
		articleId, _ := strconv.ParseInt(d.ID, 10, 64)
		authorId, _ := strconv.ParseInt(d.AuthorId, 10, 64)

		t, err := time.ParseInLocation(time.DateTime, d.PublishTime, time.Local)
		publishTimeKey := articlesKey(d.AuthorId, 0)
		likeNumKey := articlesKey(d.AuthorId, 1)

		switch int32(status) {
		case model.ArticleStatusVisible:
			l.addCacheArticleToScore(publishTimeKey, d.ID, t.Unix())
			l.addCacheArticleToScore(likeNumKey, d.ID, likeNum)
		case model.ArticleStatusUserDelete:
			_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, publishTimeKey, d.ID)
			if err != nil {
				l.Logger.Errorf("ZremCtx key: %s req: %v error: %v", publishTimeKey, d, err)
			}
			_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, likeNumKey, d.ID)
			if err != nil {
				l.Logger.Errorf("ZremCtx key: %s req: %v error: %v", likeNumKey, d, err)
			}
			// case types.ArticleStatusPending:
			// case types.ArticleStatusNotPass:
		}

		u, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
			UserId: authorId,
		})
		if err != nil {
			l.Logger.Errorf("FindById userId: %d error: %v", authorId, err)
			return err
		}
		esData = append(esData, &types.ArticleEsMsg{
			ArticleId:   articleId,
			AuthorId:    authorId,
			AuthorName:  u.Username,
			Title:       d.Title,
			Content:     d.Content,
			Description: d.Description,
			Status:      status,
			LikeNum:     likeNum,
		})
	}
	err := l.BatchUpSertToEs(l.ctx, esData)
	if err != nil {
		l.Logger.Errorf("BatchUpSertToEs data: %v error %v", esData, err)
	}
	return err
}

func (l *ArticleLogic) BatchUpSertToEs(ctx context.Context, data []*types.ArticleEsMsg) error {
	if len(data) == 0 {
		return nil
	}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "article-index",
	})
	if err != nil {
		return err
	}

	for _, d := range data {
		v, err := json.Marshal(d)
		if err != nil {
			return err
		}

		payload := fmt.Sprintf(`{"doc":%s,"doc_as_upsert":true}`, string(v))
		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "update",
			DocumentID: fmt.Sprintf("%d", d.ArticleId),
			Body:       strings.NewReader(payload),
			OnSuccess:  func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem) {},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem, err error) {
			},
		})
		if err != nil {
			return err
		}
	}
	return bi.Close(ctx)
}

func (l *ArticleLogic) addCacheArticleToScore(articleIdStr, key string, score int64) {
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if b {
		_, err := l.svcCtx.BizRedis.ZaddCtx(l.ctx, key, score, articleIdStr)
		if err != nil {
			l.Logger.Errorf("ZaddCtx key: %s error: %v", key, err)
		}
	}
}

func articlesKey(uid string, sortType int32) string {
	return fmt.Sprintf("biz#articles#%s#%d", uid, sortType)
}
