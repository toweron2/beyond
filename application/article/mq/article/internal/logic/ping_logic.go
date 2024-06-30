package logic

import (
	"context"

	"beyond/application/article/mq/article/article"
	"beyond/application/article/mq/article/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *article.Request) (*article.Response, error) {
	// todo: add your logic here and delete this line

	return &article.Response{}, nil
}
