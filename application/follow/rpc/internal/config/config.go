package config

import (
	"beyond/pkg/orm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB       orm.Config
	BizRedis redis.RedisConf
}
