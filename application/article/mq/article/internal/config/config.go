package config

import (
	"beyond/pkg/consul"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	BizRedis   redis.RedisConf
	Consul     consul.Conf
}
