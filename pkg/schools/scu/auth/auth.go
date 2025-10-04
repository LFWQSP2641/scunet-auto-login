package auth

import (
	"context"
	"fmt"
	"sync"

	C "github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/constant"
	"github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/portal/rvjx"
	S "github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/session"
)

type SCUAuthenticator struct {
	inited bool
	mu     sync.Mutex
	sess   *S.Session
}

func NewSCUAuthenticator() *SCUAuthenticator {
	return &SCUAuthenticator{
		inited: false,
		mu:     sync.Mutex{},
		sess:   &S.Session{},
	}
}

// Login 执行实际登录逻辑占位：后续可替换为真实 HTTP 请求等
func (a *SCUAuthenticator) Login(ctx context.Context, username, password string, extra map[string]string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// 这里打印以示数据已传递；生产环境应移除或使用日志
	fmt.Printf("开始登录: username=%s extra=%s\n", username, extra)

	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.inited {
		discovery := rvjx.NewDiscovery(C.MainUrl)
		preAuthInfo, err := discovery.Discover(ctx)

		if err != nil {
			return err
		}

		a.sess.PreAuth = preAuthInfo
		a.sess.Crypto = &S.CryptoContext{
			Module:            C.RSAModulus,
			PublicKeyExponent: C.RSAPublicKeyExponent,
			DeviceMAC:         preAuthInfo.DeviceMAC,
			PasswordPlain:     password,
		}

		a.inited = true
	}

	loginData := rvjx.LoginUserData{
		Username: username,
		Password: password,
		Service:  extra["service"],
	}
	loginExecutor := rvjx.NewLoginExecutor()
	err := loginExecutor.Execute(ctx, loginData, *a.sess)
	if err != nil {
		return err
	}

	return nil
}
