// Code generated by goctl. DO NOT EDIT.
// Source: follow.proto

package server

import (
	"context"

	"beyond/application/follow/rpc/internal/logic"
	"beyond/application/follow/rpc/internal/svc"
	"beyond/application/follow/rpc/pb"
)

type FollowServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedFollowServer
}

func NewFollowServer(svcCtx *svc.ServiceContext) *FollowServer {
	return &FollowServer{
		svcCtx: svcCtx,
	}
}

func (s *FollowServer) Follow(ctx context.Context, in *pb.FollowReq) (*pb.Empty, error) {
	l := logic.NewFollowLogic(ctx, s.svcCtx)
	return l.Follow(in)
}

func (s *FollowServer) UnFollow(ctx context.Context, in *pb.UnFollowReq) (*pb.Empty, error) {
	l := logic.NewUnFollowLogic(ctx, s.svcCtx)
	return l.UnFollow(in)
}

func (s *FollowServer) FollowList(ctx context.Context, in *pb.FollowListReq) (*pb.FollowListResp, error) {
	l := logic.NewFollowListLogic(ctx, s.svcCtx)
	return l.FollowList(in)
}

func (s *FollowServer) FansList(ctx context.Context, in *pb.FansListReq) (*pb.FansListResp, error) {
	l := logic.NewFansListLogic(ctx, s.svcCtx)
	return l.FansList(in)
}