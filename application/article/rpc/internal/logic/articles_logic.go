package logic

import (
	"beyond/application/article/rpc/internal/types"
	"context"
	"strconv"
	"time"

	"beyond/application/article/rpc/internal/svc"
	"beyond/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
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
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortypeInvalid
	}
	if in.UserId <= 0 {
		return nil.code.UserInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}

	sortField := "publish_time"
	var (
		sortLikeNum     int64
		sortPublishTime string
	)
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
		sortLikeNum = in.Cursor
	} else {
		sortPublishTime = time.Unix(in.Cursor, 0).Format(time.DateTime)
	}

	isCache, isEnd := false, false

	// 调用缓存查询忽略乐error,我们期望尽最大可能的给用户返回数据,不会因为redis挂掉而返回错误
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
	if len(articleIds) == int(in.PageSize) {
		// 为表示列表的结束,在Sorted Set中设置一个结束标识符,该标识符的member为-1,score为0
		// 所以我们从缓存查询出数据后,需要判断数据的最后一条是否为-1,
		// 如果为-1的话说明列表已经加载倒最后一页了,用户在滑动屏幕的话前端就不会在请求后端接口
		isCache = true
		if articleIds[len(articleIds)-1] == -1 {
			isEnd = true
		}
	}

	for _, id := range articleIds {
		l.svcCtx.ArticleModel.ArticleByUserId(l.ctx, in.UserId, in.ArticleId, in.PageSize, in.SortType)
	}

	return &pb.ArticlesResp{}, nil
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context, uid, cursor, ps int64, sortType int32) ([]int64, error) {
	key := articlesKey(uid, sortType)
	// 倒序从缓存中读取数据,并限制读条数为分页大小
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
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func articlesKey(uid int64, sortType int32) string {

}
