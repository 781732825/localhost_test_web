package model

import (
	"net/http"
	"regexp"
)

// ResponseRule 表示响应规则
type ResponseRule struct {
	Status  int               `json:"status"`  // HTTP状态码
	Headers map[string]string `json:"headers"` // 响应头
	Body    string            `json:"body"`    // 响应体
	File    string            `json:"file"`    // 本地文件路径，如果设置，将返回文件内容而不是Body
}

// Rule 表示单个匹配规则
type Rule struct {
	Path     string       `json:"path"`     // 匹配的路径（支持正则表达式）
	Method   string       `json:"method"`   // 匹配的HTTP方法（GET/POST等）
	Response ResponseRule `json:"response"` // 匹配成功后的响应

	// 编译后的正则表达式（不导出到JSON）
	pathRegex *regexp.Regexp `json:"-"`
}

// RuleFile 表示一个域名对应的规则文件
type RuleFile struct {
	Rules   []Rule       `json:"rules"`   // 规则列表
	Default ResponseRule `json:"default"` // 默认响应（当没有规则匹配时）
}

// MatchRequest 检查请求是否匹配规则
func (r *Rule) MatchRequest(req *http.Request) bool {
	// 检查方法是否匹配
	if r.Method != "" && r.Method != req.Method {
		return false
	}

	// 检查路径是否匹配
	// 如果规则的pathRegex已编译，使用正则表达式匹配
	if r.pathRegex != nil {
		return r.pathRegex.MatchString(req.URL.Path)
	}

	// 否则使用精确匹配
	return r.Path == req.URL.Path
}

// CompileRegex 编译路径中的正则表达式
func (r *Rule) CompileRegex() error {
	// 如果路径不是正则表达式，则不需要编译
	if r.Path == "" || (r.Path[0] != '^' && r.Path[len(r.Path)-1] != '$') {
		return nil
	}

	// 编译正则表达式
	regex, err := regexp.Compile(r.Path)
	if err != nil {
		return err
	}

	// 保存编译后的正则表达式
	r.pathRegex = regex
	return nil
}