package rvjx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	C "github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/constant"
	scuerror "github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/error"
)

func Logout(ctx context.Context) error {
	logoutURL := C.LogoutPostUrl
	req, err := http.NewRequestWithContext(ctx, "POST", logoutURL, nil)
	if err != nil {
		return err
	}

	uaHeader := C.UserAgent
	req.Header.Set("User-Agent", uaHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return scuerror.ErrConnection
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return parseLogoutResponse(string(body))
}

func parseLogoutResponse(responseText string) error {
	if strings.Contains(responseText, `"result":"success"`) {
		return nil
	}

	// 登录失败，提取错误信息
	errorMsg := ExtractErrorMessage(responseText)

	// 根据错误信息返回特定异常
	if strings.Contains(errorMsg, "用户已不在线") {
		return scuerror.ErrUserNotLogin
	}

	return fmt.Errorf("%w: %s", scuerror.ErrLogoutFailed, errorMsg)
}
