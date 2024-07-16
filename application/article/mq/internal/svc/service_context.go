package svc

import (
	"beyond/application/article/internal/model"
	"beyond/application/article/mq/internal/config"
	"beyond/application/user/rpc/user"
	"beyond/pkg/es"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis

	UserRpc user.User
	Es      *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdf := redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	}
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(conn, cache.CacheConf{{rdf, 100}}),
		BizRedis:     redis.MustNewRedis(rdf),
		UserRpc:      user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		Es:           es.MustNewEs(&c.Es),
	}
}
