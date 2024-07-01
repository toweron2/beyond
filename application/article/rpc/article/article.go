// Code generated by goctl. DO NOT EDIT.
// Source: article.proto

package article

import (
	"context"

	"beyond/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ArticleDeleteReq  = pb.ArticleDeleteReq
	ArticleDetailReq  = pb.ArticleDetailReq
	ArticleDetailResp = pb.ArticleDetailResp
	ArticleItem       = pb.ArticleItem
	ArticlesReq       = pb.ArticlesReq
	ArticlesResp      = pb.ArticlesResp
	Empty             = pb.Empty
	PublishReq        = pb.PublishReq
	PublishResp       = pb.PublishResp

	Article interface {
		Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error)
		Articles(ctx context.Context, in *ArticlesReq, opts ...grpc.CallOption) (*ArticlesResp, error)
		ArticleDelete(ctx context.Context, in *ArticleDeleteReq, opts ...grpc.CallOption) (*Empty, error)
		ArticleDetail(ctx context.Context, in *ArticleDetailReq, opts ...grpc.CallOption) (*ArticleDetailResp, error)
	}

	defaultArticle struct {
		cli zrpc.Client
	}
)

func NewArticle(cli zrpc.Client) Article {
	return &defaultArticle{
		cli: cli,
	}
}

func (m *defaultArticle) Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishResp, error) {
	client := pb.NewArticleClient(m.cli.Conn())
	return client.Publish(ctx, in, opts...)
}

func (m *defaultArticle) Articles(ctx context.Context, in *ArticlesReq, opts ...grpc.CallOption) (*ArticlesResp, error) {
	client := pb.NewArticleClient(m.cli.Conn())
	return client.Articles(ctx, in, opts...)
}

func (m *defaultArticle) ArticleDelete(ctx context.Context, in *ArticleDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	client := pb.NewArticleClient(m.cli.Conn())
	return client.ArticleDelete(ctx, in, opts...)
}

// rpc ArticleDelete(ArticleDeleteReq) returns (google.protobuf.Empty);
func (m *defaultArticle) ArticleDetail(ctx context.Context, in *ArticleDetailReq, opts ...grpc.CallOption) (*ArticleDetailResp, error) {
	client := pb.NewArticleClient(m.cli.Conn())
	return client.ArticleDetail(ctx, in, opts...)
}
