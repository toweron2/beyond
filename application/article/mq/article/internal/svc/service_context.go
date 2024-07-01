package svc

import (
	"beyond/application/article/mq/article/internal/config"
	"beyond/application/user/rpc/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config   config.Config
	BizRedis *redis.Redis
	UserRpc  user.User
	Es       *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdf := redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	}

	return &ServiceContext{
		Config:   c,
		BizRedis: redis.MustNewRedis(rdf),
	}
}
