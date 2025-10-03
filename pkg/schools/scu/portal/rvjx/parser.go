package rvjx

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// ParseQueryString 解析查询字符串，提取参数
func ParseQueryString(queryString string) (map[string]string, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for key, valueList := range values {
		if len(valueList) > 0 {
			result[key] = valueList[0]
		}
	}

	return result, nil
}

// ExtractQueryStringFromHTML 从HTML中提取查询字符串
func ExtractQueryStringFromHTML(html string) (string, error) {
	// 尝试正则匹配
	re := regexp.MustCompile(`/index\.jsp\?([^\'\"]+)`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1], nil
	}

	// 回退到字符串查找方法
	startIndex := strings.Index(html, "/index.jsp?")
	if startIndex == -1 {
		return "", errors.New("无法找到查询字符串")
	}

	startIndex += 11 // len("/index.jsp?")
	endIndex := strings.Index(html[startIndex:], "'</script>")
	if endIndex == -1 {
		return "", errors.New("无法确定查询字符串结束位置")
	}

	return html[startIndex : startIndex+endIndex], nil
}
