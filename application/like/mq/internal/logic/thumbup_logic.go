package logic

import (
	"beyond/application/like/mq/internal/svc"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "mq like", Value: "thumbup_logic"}),
	}
}

// 消费kafka消息
func (l *ThumbupLogic) Consume(key, val string) error {
	fmt.Printf("get key: %s val: %s\n", key, val)
	return nil
}
