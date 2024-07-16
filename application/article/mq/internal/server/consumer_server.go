package server

import (
	"beyond/application/article/mq/internal/logic"
	"beyond/application/article/mq/internal/svc"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, logic.NewArticleLikeNumLogic(ctx, svcCtx)),
		kq.MustNewQueue(svcCtx.Config.ArticleKqConsumerConf, logic.NewArticleLogic(ctx, svcCtx)),
	}
}
