package server

import (
	"beyond/application/like/mq/internal/logic"
	"beyond/application/like/mq/internal/svc"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, logic.NewThumbupLogic(ctx, svcCtx)),
	}
}
