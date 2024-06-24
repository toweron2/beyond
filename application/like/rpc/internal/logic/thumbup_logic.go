package logic

import (
	"beyond/application/like/rpc/internal/types"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/threading"

	"beyond/application/like/rpc/internal/svc"
	"beyond/application/like/rpc/service"

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
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "rpc like", Value: "thumbup_logic"}),
	}
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupReq) (*service.ThumbupResp, error) {
	// TODO 逻辑暂时忽略
	// 1. 查询是否点过赞
	// 2. 计算当前内容的总点赞数和点踩数
	msg := &types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}

	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Thumbup] marshal msg: %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(string(data))
		if err != nil {
			l.Logger.Errorf("[Thumbup] kq push data: %s error: %v", data, err)
		}
	})

	return &service.ThumbupResp{}, nil
}
