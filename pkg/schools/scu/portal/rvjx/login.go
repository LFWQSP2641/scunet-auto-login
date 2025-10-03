package rvjx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	C "scunet-auto-login/pkg/schools/scu/constant"
	scucry "scunet-auto-login/pkg/schools/scu/crypto"
	scuerror "scunet-auto-login/pkg/schools/scu/error"
	S "scunet-auto-login/pkg/schools/scu/session"
	"strings"
)

type LoginUserData struct {
	Username string
	Password string
	Service  string
}

// LoginExecutor 登录执行器
type LoginExecutor struct {
}

// NewLoginExecutor 创建登录执行器
func NewLoginExecutor() *LoginExecutor {
	return &LoginExecutor{}
}

// Execute 执行登录操作
func (le *LoginExecutor) Execute(ctx context.Context, user LoginUserData, session S.Session) error {
	// 验证服务类型
	serviceCode, exists := C.ServiceCodes()[user.Service]
	if !exists {
		return fmt.Errorf("%w: %s", scuerror.ErrLoginService, user.Service)
	}

	// 加密密码
	encryptedPassword, err := scucry.EncryptedPassword(*session.Crypto)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 构建登录数据
	loginData := url.Values{
		"userId":          {user.Username},
		"password":        {encryptedPassword},
		"service":         {serviceCode},
		"queryString":     {session.PreAuth.LoginQuery},
		"operatorPwd":     {""},
		"operatorUserId":  {""},
		"validcode":       {""},
		"passwordEncrypt": {"true"},
	}

	// 清除cookies避免风控
	// 这里需要根据实际的HTTP客户端实现来清除cookies

	// 构建请求
	loginURL := C.LoginPostUrl
	req, err := http.NewRequestWithContext(ctx, "POST", loginURL, strings.NewReader(loginData.Encode()))
	if err != nil {
		return err
	}

	// 设置请求头
	for key, value := range C.HttpHeader() {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return scuerror.ErrLoginConnection
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析响应
	return le.parseLoginResponse(string(body))
}

// parseLoginResponse 解析登录响应
func (le *LoginExecutor) parseLoginResponse(responseText string) error {
	// 处理Unicode编码
	responseText = strings.ReplaceAll(responseText, "\\u", "\\u")

	if strings.Contains(responseText, `"result":"success"`) {
		return nil
	}

	// 登录失败，提取错误信息
	errorMsg := le.extractErrorMessage(responseText)

	// 根据错误信息返回特定异常
	if strings.Contains(errorMsg, "在线用户数量上限") {
		return scuerror.ErrOnlineUserLimit
	}
	if strings.Contains(errorMsg, "验证码错误") {
		return scuerror.ErrLoginRiskControl
	}

	return fmt.Errorf("%w: %s", scuerror.ErrLoginFailed, errorMsg)
}

// extractUserIndex 提取用户索引
func (le *LoginExecutor) extractUserIndex(responseText string) string {
	// 使用正则表达式提取
	re := regexp.MustCompile(`userIndex":"([^"]+)"`)
	match := re.FindStringSubmatch(responseText)
	if len(match) > 1 {
		return match[1]
	}

	// 回退方法
	startIndex := strings.Index(responseText, "userIndex") + 12
	endIndex := strings.Index(responseText, `","result"`)
	if startIndex > 11 && endIndex > startIndex {
		return responseText[startIndex:endIndex]
	}

	return ""
}

// extractErrorMessage 提取错误信息
func (le *LoginExecutor) extractErrorMessage(responseText string) string {
	re := regexp.MustCompile(`"message":"([^"]+)"`)
	match := re.FindStringSubmatch(responseText)
	if len(match) > 1 {
		return match[1]
	}

	// 回退方法
	startIndex := strings.Index(responseText, `"message"`) + 11
	endIndex := strings.Index(responseText, `","forwordurl"`)
	if startIndex > 10 && endIndex > startIndex {
		return responseText[startIndex:endIndex]
	}

	return "未知错误"
}

//// Logout 执行登出操作
//func (le *LoginExecutor) Logout(ctx context.Context, session *auth.Session) error {
//	if session == nil || session.UserIndex == "" {
//		return auth.ErrInvalidSession
//	}
//
//	logoutURL := le.portal.baseURL + "eportal/InterFace.do?method=logout"
//
//	req, err := http.NewRequestWithContext(ctx, "POST", logoutURL, nil)
//	if err != nil {
//		return err
//	}
//
//	// 设置请求头
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
//	req.Header.Set("User-Agent", le.portal.GetOptions().UserAgent)
//	req.Header.Set("Accept", "*/*")
//	req.Header.Set("Accept-Encoding", "gzip, deflate")
//
//	// 发送请求
//	client := le.portal.GetClient()
//	resp, err := client.Do(req)
//	if err != nil {
//		return auth.ErrLoginConnection
//	}
//	defer resp.Body.Close()
//
//	// 检查响应状态
//	if resp.StatusCode != http.StatusOK {
//		return auth.ErrLogoutFailed
//	}
//
//	return nil
//}
