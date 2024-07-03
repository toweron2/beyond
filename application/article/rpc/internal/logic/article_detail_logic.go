package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"beyond/application/article/rpc/internal/svc"
	"beyond/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "rpc article", Value: "article_detail_logic"}),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *pb.ArticleDetailReq) (*pb.ArticleDetailResp, error) {
	// go-zero框架默认会给不存在的加上空缓存
	// get cache:beyondArticle:article:id:999
	// "*"
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return &pb.ArticleDetailResp{}, nil
		}
		return nil, err
	}
	return &pb.ArticleDetailResp{
		Article: &pb.ArticleItem{
			Id:          article.Id,
			Title:       article.Title,
			Content:     article.Content,
			Description: article.Description,
			Cover:       article.Cover,
			AuthorId:    article.AuthorId,
			LikeCount:   article.LikeNum,
			PublishTime: article.PublishTime.Unix(),
		},
	}, nil
}