package logic

import (
	"beyond/application/applet/internal/code"
	"beyond/application/user/rpc/user"
	"beyond/pkg/encrypt"
	"beyond/pkg/jwt"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strings"

	"beyond/application/applet/internal/svc"
	"beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const prefixActivation = "biz#activation#count#%s"

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	if req.Name = strings.TrimSpace(req.Name); len(req.Name) == 0 {
		l.Logger.Error("name cannot be empty")
		return nil, code.RegisterNameEmpty
	}
	if req.Mobile = strings.TrimSpace(req.Mobile); len(req.Name) == 0 {
		l.Logger.Error("mobile cannot be empty")
		return nil, code.RegisterMobileEmpty
	}
	if req.Password = strings.TrimSpace(req.Password); len(req.Name) == 0 {
		l.Logger.Error("password cannot be empty")
		return nil, code.RegisterPasswdEmpty
	}
	req.Password = encrypt.EncPassword(req.Password)

	if req.VerificationCode = strings.TrimSpace(req.VerificationCode); len(req.VerificationCode) == 0 {
		l.Logger.Error("verification code cannot be empty")
		return nil, code.VerificationCodeEmpty
	}
	err = checkVerificationCode(req.Mobile, req.VerificationCode, l.svcCtx.BizRedis)
	if err != nil {
		l.Logger.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		l.Logger.Errorf("EncMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	userRet, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		l.Logger.Errorf("FindByMobile error: %v", err)
		return nil, err
	}
	if userRet != nil && userRet.UserId > 0 {
		return nil, errors.New("this mobile is already registered")
	}
	regRet, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   mobile,
	})
	if err != nil {
		return nil, err
	}
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]any{
			"userId": regRet.UserId,
		},
	})
	if err != nil {
		l.Logger.Errorf("BuildTokens error: %v", err)
		return nil, err
	}
	err = delActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		l.Logger.Errorf("delActivationCache error: %v", err)
	}

	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, err
}

func checkVerificationCode(mobile, code string, rds *redis.Redis) error {
	cacheCode, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return errors.New("verification  code expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}
	return nil
}
