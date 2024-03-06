package logic

import (
    "beyond/application/user/rpc/internal/svc"
    "beyond/application/user/rpc/service"
    "context"

    "github.com/zeromicro/go-zero/core/logx"
)

type FindByMobileLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
    return &FindByMobileLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *FindByMobileLogic) FindByMobile(in *service.FindByMobileRequest) (*service.FindByMobileResponse, error) {
    // TODO find mobile

    return &service.FindByMobileResponse{}, nil
}