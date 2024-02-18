package svc

import (
	"beyond/application/applet/internal/config"
	"beyond/application/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRpc  user.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 自定义拦截器
	// userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor*()))
	rds, err := redis.NewRedis(c.BizRedis)
	if err != nil {
		logx.Errorf("errors : %v", err)
		return nil
	}
	return &ServiceContext{
		Config:   c,
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		BizRedis: rds,
	}
}
