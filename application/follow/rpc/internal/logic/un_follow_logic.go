package logic

import (
	"context"

	"beyond/application/follow/rpc/internal/svc"
	"beyond/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "rpc follow", Value: "un_follow_logic"}),
	}
}

func (l *UnFollowLogic) UnFollow(in *pb.UnFollowReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line
	if in.UserId <= 0 {
		// return n
	}
	return &pb.Empty{}, nil
}
