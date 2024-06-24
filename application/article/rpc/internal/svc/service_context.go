package svc

import (
	"beyond/application/article/rpc/internal/config"
	"beyond/application/article/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdf := redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	}

	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(sqlx.NewMysql(c.DataSource), cache.CacheConf{{rdf, 100}}),
		BizRedis:     redis.MustNewRedis(rdf),
	}

}
