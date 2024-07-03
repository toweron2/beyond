package svc

import (
	"beyond/application/article/internal/model"
	"beyond/application/article/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	ArticleModel      model.ArticleModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config: c,
		ArticleModel: model.NewArticleModel(
			sqlx.NewMysql(c.DataSource),
			cache.CacheConf{{c.BizRedis, 100}}),
		BizRedis: redis.MustNewRedis(c.BizRedis),
	}

}
