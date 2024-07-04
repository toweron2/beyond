package logic

import (
	"beyond/application/follow/rpc/internal/code"
	"beyond/application/follow/rpc/internal/model"
	"context"
	"gorm.io/gorm"
	"strconv"

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
	if in.UserId <= 0 {
		return nil, code.UserIdEmpty
	}
	if in.FollowedUserId <= 0 {
		return nil, code.FollowedUserIdEmpyt
	}
	follow, err := l.svcCtx.FollowModel.FindByUserIDAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[UnFollow] FollowModel.FindByUserIDAndFollowedUserID err: %v req: %v", err, in)
		return nil, err
	}
	if follow == nil {
		return nil, nil
	}

	// 事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]any{
			"follow_status": model.FollowStatusUnfollow,
		})
		if err != nil {
			return err
		}
		err = model.NewFollowCountModel(tx).DecrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		return model.NewFollowCountModel(tx).DecrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		l.Logger.Errorf("[UnFollow] Transtction error: %v", err)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFollowKey(in.UserId), strconv.FormatInt(in.FollowedUserId, 10))
	if err != nil {
		l.Logger.Errorf("[UnFollow] Redis Zrem error: %v", err)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFansKey(in.FollowedUserId), strconv.FormatInt(in.UserId, 10))
	if err != nil {
		l.Logger.Errorf("[UnFollow] Redis Zrem error: %v", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}
