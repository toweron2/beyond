package svc

import (
	"beyond/application/article/api/internal/config"
	"beyond/application/article/rpc/article"
	"beyond/pkg/oss"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	Oss        oss.OSS
	ArticleRPC article.Article
}

func NewServiceContext(c config.Config) *ServiceContext {
	qiniu := oss.NewQiniu(c.Qiniu.AccessKey, c.Qiniu.SecretKey, c.Qiniu.Bucket, c.Qiniu.UseHTTPS, c.Qiniu.UseCdnDomains)
	return &ServiceContext{
		Config:     c,
		Oss:        qiniu,
		ArticleRPC: article.NewArticle(zrpc.MustNewClient(c.ArticleRPC)),
	}
}
