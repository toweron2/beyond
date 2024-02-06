package logic

import (
	"context"

	"beyond/application/applet/internal/svc"
	"beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

//	             _ooOoo_
//	            o8888888o
//	            88" . "88
//	            (| -_- |)
//	             O\ = /O
//	         ____/`---'\____
//	       .   ' \\| |// `.
//	        / \\||| : |||// \
//	      / _||||| -:- |||||- \
//	        | | \\\ - /// | |
//	      | \_| ''\---/'' | |
//	       \ .-\__ `-` ___/-. /
//	    ___`. .' /--.--\ `. . __
//	 ."" '< `.___\_<|>_/___.' >'"".
//	| | : `- \`.;`\ _ /`;.`/ - ` : | |
//	  \ \ `-. \_ __\ /__ _/ .-` / /
//
// ======`-.____`-.___\_____/___.-`____.-'======
//
//	`=---='
//
// .............................................
//
//	佛祖保佑             永无BUG
type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	return
}
