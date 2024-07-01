package logic

import (
	"beyond/application/article/mq/article/internal/svc"
	"beyond/application/article/mq/article/internal/types"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type ArticleLikeNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLikeNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLikeNumLogic {
	return &ArticleLikeNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx).WithFields(logx.LogField{Key: "mq article", Value: "article_likenum_logic"}),
	}
}

func (l *ArticleLikeNumLogic) Consume(_, val string) error {
	var msg *types.CanalLikeMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		l.Logger.Errorf("Consume val: %s error: %v", val, err)
		return err
	}
	return l.updateArticleLikeNum(l.ctx, msg)
}

func (l *ArticleLikeNumLogic) updateArticleLikeNum(ctx context.Context, msg *types.CanalLikeMsg) error {
	if len(msg.Data) == 0 {
		for _, d := range msg.Data {
			if d.BizID != types.ArticleBizID {
				continue
			}
			id, err := strconv.ParseInt(d.ObjID, 10, 64)
			if err != nil {
				l.Logger.Errorf("strconv.ParseInt id: %s error: %v", d.ID, err)
				continue
			}
			likeNum, err := strconv.ParseInt(d.LikeNum, 10, 64)
			if err != nil {
				l.Logger.Errorf("strconv.ParseInt likeNum: %s error: %v", d.LikeNum, err)
				continue
			}
			err = l.svcCtx.ArticleModel.UpdateLikeNum(ctx, id, likeNum)
			if err != nil {
				logx.Errorf("UpdateLikeNum id: %d like: %d", id, likeNum)
			}
		}
	}
	return nil
}
