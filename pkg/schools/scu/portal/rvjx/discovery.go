package rvjx

import (
	"context"
	"io"
	"net/http"
	scuerror "scunet-auto-login/pkg/schools/scu/error"
	S "scunet-auto-login/pkg/schools/scu/session"
)

// Discovery 预认证发现器 (相当于Python中的getQueryString)
type Discovery struct {
	baseUrl string
}

// NewDiscovery 创建发现器实例
func NewDiscovery(baseUrl string) *Discovery {
	return &Discovery{
		baseUrl: baseUrl,
	}
}

// GetQueryString 获取登录所需的查询字符串参数
func (d *Discovery) GetQueryString(ctx context.Context) (string, error) {
	// 发起初始请求获取重定向
	req, err := http.NewRequestWithContext(ctx, "GET", d.baseUrl, nil)
	if err != nil {
		return "", err
	}

	// 设置不自动跟随重定向
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", scuerror.ErrLoginConnection
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", scuerror.ErrLoginConnection
	}

	bodyString := string(bodyBytes)

	// 检查是否已经登录
	if resp.Request.URL.Path == "/eportal/success.jsp" {
		return "", scuerror.ErrAlreadyLoggedIn
	}

	// 提取查询字符串
	queryString, err := ExtractQueryStringFromHTML(bodyString)
	if err != nil {
		return "", scuerror.ErrLoginParameter
	}

	// 缓存结果
	return queryString, nil
}

// Discover 实现PreAuthenticator接口
func (d *Discovery) Discover(ctx context.Context) (*S.PreAuthInfo, error) {
	queryString, err := d.GetQueryString(ctx)
	if err != nil {
		return nil, err
	}

	// 解析查询字符串获取参数
	params, err := ParseQueryString(queryString)
	if err != nil {
		return nil, scuerror.ErrLoginParameter
	}

	return &S.PreAuthInfo{
		LoginQuery:       queryString,
		DeviceMAC:        params["mac"],
		AdditionalParams: params,
	}, nil
}
