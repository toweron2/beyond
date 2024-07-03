package logic

import (
	"context"

	"beyond/application/follow/rpc/internal/svc"
	"beyond/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *pb.FollowListReq) (*pb.FollowListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.FollowListResp{}, nil
}
