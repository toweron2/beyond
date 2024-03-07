package main

import (
	"beyond/application/user/rpc/internal/config"
	"beyond/application/user/rpc/internal/server"
	"beyond/application/user/rpc/internal/svc"
	"beyond/application/user/rpc/service"
	"beyond/pkg/interceptors"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	cs "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		service.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == cs.DevMode || c.Mode == cs.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 自定义拦截器
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())

	defer s.Stop()

	// 服务注册
	// consul.Re

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
