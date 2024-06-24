package logic

import (
	"beyond/application/user/rpc/internal/code"
	"beyond/application/user/rpc/internal/model"
	"beyond/application/user/rpc/internal/svc"
	"beyond/application/user/rpc/service"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "rpc  user", Value: "register_logic"}),
	}
}

func (l *RegisterLogic) Register(in *service.RegisterRequest) (*service.RegisterResponse, error) {
	// 的注册名抽为空的是脏,返回业务自定错误码
	if len(in.Username) == 0 {
		return nil, code.RegisterNameEmpty
	}
	now := time.Now()
	ret, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username:   in.Username,
		Mobile:     in.Mobile,
		Avatar:     in.Avatar,
		CreateTime: now,
		UpdateTime: now,
	})
	if err != nil {
		l.Logger.Errorf("Register req: %v error: %v", in, err)
		return nil, err
	}
	userId, err := ret.LastInsertId()
	if err != nil {
		l.Logger.Errorf("Register LastInsertId error: %v", err)
		return nil, err
	}

	return &service.RegisterResponse{UserId: userId}, nil
}
