package error

import "errors"

var (
	ErrLoginConnection  = errors.New("网络连接问题")
	ErrLoginRedirect    = errors.New("重定向异常")
	ErrLoginParameter   = errors.New("参数提取异常")
	ErrLoginService     = errors.New("服务类型异常")
	ErrLoginProcess     = errors.New("登录过程异常")
	ErrAlreadyLoggedIn  = errors.New("已经登录")
	ErrLoginFailed      = errors.New("登录失败")
	ErrOnlineUserLimit  = errors.New("在线用户数量上限")
	ErrLoginRiskControl = errors.New("登录风控")
)
