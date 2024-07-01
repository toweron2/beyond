package logic

import (
	"beyond/application/article/model"
	"beyond/application/article/rpc/internal/code"
	"beyond/application/article/rpc/internal/svc"
	"beyond/application/article/rpc/internal/types"
	"beyond/application/article/rpc/pb"
	"context"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "rpc article", Value: "publish_logic"}),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishReq) (*pb.PublishResp, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.Title == "" {
		return nil, code.ArticleTitleCantEmpty
	}
	if in.Content == "" {
		return nil, code.ArticleContentCantEmpty
	}

	ret, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		AuthorId:    in.UserId,
		Title:       in.Title,
		Content:     in.Content,
		Description: in.Description,
		Cover:       in.Cover,
		Status:      int64(model.ArticleStatusVisible), // 正常逻辑不会这样写，这里为了演示方便
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		l.Logger.Errorf("Publish Insert req: %v error: %v", in, err)
		return nil, err
	}

	articleId, err := ret.LastInsertId()
	if err != nil {
		l.Logger.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	var (
		articleIdStr   = strconv.FormatInt(articleId, 10)
		publishTimeKey = articlesKey(in.UserId, types.SortPublishTime)
		likeNumKey     = articlesKey(in.UserId, types.SortLikeCount)
	)
	l.addCacheArticleToScore(articleIdStr, publishTimeKey, time.Now().Unix())
	l.addCacheArticleToScore(articleIdStr, likeNumKey, 0)

	return &pb.PublishResp{ArticleId: articleId}, nil
}

func (l *PublishLogic) addCacheArticleToScore(articleIdStr, key string, score int64) {
	// 需要判断是否存在缓存
	/*
		不判断是用户是否存在文章缓存, 没没有缓存的的e情况下,只将最新发布的添加倒缓存中了
		当下次进行查询的时候,缓存存在,但是内容只有上一次新发布的,导致查出来数据不一致
	*/
	b, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, key)
	if b {
		_, err := l.svcCtx.BizRedis.ZaddCtx(l.ctx, key, score, articleIdStr)
		if err != nil {
			logx.Errorf("ZaddCtx key: %s error: %v", key, err)
		}
	}
}
