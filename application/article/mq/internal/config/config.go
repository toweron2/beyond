package config

import (
	"beyond/pkg/es"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf

	KqConsumerConf        kq.KqConf
	ArticleKqConsumerConf kq.KqConf
	DataSource            string
	BizRedis              redis.RedisConf
	// Consul                consul.Conf
	Es      es.Config
	UserRPC zrpc.RpcClientConf
}
