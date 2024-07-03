package main

import (
	"beyond/application/article/mq/internal/config"
	"beyond/application/article/mq/internal/svc"
	"context"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.ServiceConf.MustSetUp()

	ctx := svc.NewServiceContext(c)

	s := service.NewServiceGroup()
	defer s.Stop()

	for _, mq := range svc.Consumers(context.Background(), ctx) {
		s.Add(mq)
	}

	fmt.Printf("Starting rpc server at %s...\n")
	s.Start()
}
