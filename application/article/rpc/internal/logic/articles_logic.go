package logic

import (
	"beyond/application/article/rpc/internal/code"
	"beyond/application/article/rpc/internal/model"
	"beyond/application/article/rpc/internal/types"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"sort"
	"strconv"
	"time"

	"beyond/application/article/rpc/internal/svc"
	"beyond/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixArticles = "biz#articles#%d#%d"
	articlesExpire = 3600 * 24 * 2
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "rpc article", Value: "articles_logic"}),
	}
}

/*
Articles
查看某个用户发布过的文章需要支持分页,通过往上滑动可以不断地加载下一页,
文章支持按照点赞数和发布时间倒序返回列表,使用游标地方式进行分页.
使用sorted set 来存储,member为文章地id,即我们在Sorted Set中存储缓存索引,
查出来缓存索引后,因为我们自动生成乐以主键id索引为key的缓存,
所以查出来索引列表后我们在查询行记录缓存极客获取文章详情,
Sorted Set的score为文章的点赞数或者文章发布时间.
*/
func (l *ArticlesLogic) Articles(in *pb.ArticlesReq) (*pb.ArticlesResp, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = types.DefaultSortLikeCursor
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		}
	}

	var (
		sortField       string
		sortLikeNum     int64
		sortPublishTime string
	)
	switch in.SortType {
	case types.SortPublishTime:
		sortField = "publish_time"
		sortPublishTime = time.Unix(in.Cursor, 0).Format(time.DateTime)
	case types.SortLikeCount:
		sortField = "like_num"
		sortLikeNum = in.Cursor
	}

	isEnd := false
	var retArticles []*model.Article

	// 调用缓存查询忽略了error,我们期望尽最大可能的给用户返回数据,不会因为redis挂掉而返回错误
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
	if len(articleIds) > 0 {
		if articleIds[len(articleIds)-1] == -1 {
			// 为表示列表的结束,在Sorted Set中设置一个结束标识符,该标识符的member为-1,score为0
			// 所以我们从缓存查询出数据后,需要判断数据的最后一条是否为-1,
			// 如果为-1的话说明列表已经加载倒最后一页了,用户在滑动屏幕的话前端就不会在请求后端接口
			isEnd = true
		}
		articles, err := l.articleByIds(l.ctx, articleIds)
		if err != nil {
			return nil, err
		}

		// 通过sortFiled对artcles进行排序 go 1.21
		/*var cmpFunc func(a, b *model.Article) int
		  if sortField == "like_num" {
		  	cmpFunc = func(a, b *model.Article) int {
		  		return cmp.Compare(b.LikeNum, a.LikeNum)
		  	}
		  } else {
		  	cmpFunc = func(a, b *model.Article) int {
		  		return cmp.Compare(b.PublishTime.Unix(), a.PublishTime.Unix())
		  	}
		}*/
		var cmpFunc func(i, j int) bool
		if sortField == "like_num" {
			cmpFunc = func(i, j int) bool {
				return articles[i].LikeNum > articles[j].LikeNum
			}
		} else {
			cmpFunc = func(i, j int) bool {
				return articles[i].PublishTime.Unix() > articles[j].PublishTime.Unix()
			}
		}
		sort.Slice(articles, cmpFunc)
		retArticles = articles
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("ArticleByUserId:%d:%d", in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.ArticleModel.ArticlesByUserId(l.ctx, in.UserId, sortLikeNum, types.DefaultLimit, types.ArticleStatusVisible, sortPublishTime, sortField)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d sortField: %s error: %v", in.UserId, sortField, err)
			return nil, err
		}
		if v == nil {
			return &pb.ArticlesResp{}, nil
		}
		articles := v.([]*model.Article)
		if len(articles) > int(in.PageSize) {
			retArticles = articles[:int(in.PageSize)]
		} else {
			retArticles = articles
			isEnd = true
		}

		threading.GoSafe(func() {
			// 异步添加到缓存
			if len(articles) > 0 && len(articles) < types.DefaultLimit {
				articles = append(articles, &model.Article{Id: -1})
			}
			err = l.addCacheArticles(context.Background(), articles, in.UserId, in.SortType)
			if err != nil {
				logx.Errorf("addCacheArticles error: %v", err)
			}
		})
	}

	curPage := make([]*pb.ArticleItem, 0, len(retArticles))
	for _, article := range retArticles {
		curPage = append(curPage, &pb.ArticleItem{
			Id:           article.Id,
			Title:        article.Title,
			Content:      article.Content,
			Description:  article.Description, // nil
			Cover:        article.Cover,       // nil
			CommentCount: article.CommentNum,
			LikeCount:    article.LikeNum,
			PublishTime:  article.PublishTime.Unix(),
		})
	}

	var lastId, cursor int64
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		switch in.SortType {
		case types.SortPublishTime:
			cursor = pageLast.PublishTime
		case types.SortLikeCount:
			cursor = pageLast.LikeCount
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, article := range curPage {
			if article.Id == in.ArticleId {
				if (in.SortType == types.SortPublishTime && article.PublishTime == in.Cursor) ||
					(in.SortType == types.SortLikeCount && article.LikeCount == in.Cursor) {
					curPage = curPage[k:]
					break
				}
			}
		}
	}

	return &pb.ArticlesResp{
		Articles:  curPage,
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
	}, nil
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context, uid, cursor, ps int64, sortType int32) ([]int64, error) {
	key := articlesKey(uid, sortType)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}

	// 按分数倒序从缓存中读取数据,并限制读条数为分页大小
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(ps))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimitCtx key: %s error: %v", key, err)
		return nil, err
	}
	ids := make([]int64, 0, len(pairs))
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (l *ArticlesLogic) articleByIds(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](func(source chan<- int64) {
		for _, aid := range articleIds {
			if aid == -1 {
				continue
			}
			source <- aid
		}
	}, func(id int64, writer mr.Writer[*model.Article], cancel func(error)) {
		p, err := l.svcCtx.ArticleModel.FindOne(ctx, id)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(p)
	}, func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) {
		articles := make([]*model.Article, 0)
		for article := range pipe {
			articles = append(articles, article)
		}
		writer.Write(articles)
	})
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (l *ArticlesLogic) addCacheArticles(ctx context.Context, articles []*model.Article, uid int64, sortType int32) error {
	if len(articles) == 0 {
		return nil
	}
	key := articlesKey(uid, sortType)
	for _, article := range articles {
		score := article.LikeNum
		if sortType == types.SortPublishTime {
			score = article.PublishTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(article.Id, 64))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
}

func articlesKey(uid int64, sortType int32) string {
	return fmt.Sprintf(prefixArticles, uid, sortType)
}
