package svc

import (
	"beyond/application/applet/internal/config"
	"beyond/application/user/rpc/user"
	"beyond/pkg/interceptors"
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
	clientErrorInterceptor := zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor())
	// userRPC := zrpc.MustNewClient(c.UserRPC, clientErrorInterceptor)
	rds, err := redis.NewRedis(c.BizRedis)
	if err != nil {
		logx.Errorf("errors : %v", err)
		return nil
	}
	return &ServiceContext{
		Config:   c,
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRPC, clientErrorInterceptor)),
		BizRedis: rds,
	}
}
