package logic

import (
	"beyond/application/user/user"
	"beyond/pkg/util"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	"beyond/application/applet/internal/svc"
	"beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	verificationLimitCountKey  = ""
	prefixVerificationCountKey = "biz#verification#count#%s"
	verificationLimitPerDay    = 10
	expireActivation           = 1800
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		return nil, err
	}

	// 30分钟内验证码不再变化
	code, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache mobile: %s error: %v", req.Mobile, err)
	}
	if len(code) == 0 {
		code = util.RandomNumeric(6)
	}
	_, err = l.svcCtx.UserRpc.SendSms(l.ctx, &user.SendSmsRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
	}
	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}

	return &types.VerificationResponse{}, nil
}

func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCountKey, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

// incrVerificationCount 增加验证码计数并设置过期时间, 该方法通过给定的手机号码，使用 Redis 存储验证码计数，并设置计数键的过期时间。
func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCountKey, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}
	return l.svcCtx.BizRedis.Expireat(key, time.Now().Truncate(time.Hour*24).Add(time.Hour*24-1).Unix())
}

func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	return rds.Get(fmt.Sprintf(prefixActivation, mobile))
}

func saveActivationCache(mobile, code string, rds *redis.Redis) error {
	return rds.Setex(fmt.Sprintf(prefixActivation, mobile), code, expireActivation)
}

func delActivationCache(mobile string, rds *redis.Redis) error {
	_, err := rds.Del(fmt.Sprintf(prefixActivation, mobile))
	return err
}
