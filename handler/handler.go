package handler

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// RequestHandler 处理HTTP请求
type RequestHandler struct {
	ruleManager *RuleManager
}

// NewRequestHandler 创建请求处理器
func NewRequestHandler(ruleManager *RuleManager) *RequestHandler {
	return &RequestHandler{
		ruleManager: ruleManager,
	}
}

// ServeHTTP 实现http.Handler接口
func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 记录请求信息
	log.Printf("收到请求: %s %s %s", r.Host, r.Method, r.URL.Path)

	// 记录查询参数
	if len(r.URL.RawQuery) > 0 {
		log.Printf("请求参数: %s", r.URL.RawQuery)
	}

	// 提取域名（用于调试）
	domain := r.Host
	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}
	log.Printf("处理域名: %s", domain)

	// 匹配规则
	responseRule, err := h.ruleManager.MatchRule(r)
	if err != nil {
		// 处理错误
		log.Printf("匹配规则失败: %v", err)
		// 设置默认的Content-Type和UTF-8编码
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("服务器内部错误"))
		return
	}

	// 设置响应头
	for key, value := range responseRule.Headers {
		w.Header().Set(key, value)
	}

	// 确保Content-Type包含charset=utf-8
	contentType := w.Header().Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "charset=") {
		w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	}

	// 设置状态码
	w.WriteHeader(responseRule.Status)

	// 检查是否需要返回文件内容
	if responseRule.File != "" {
		// 读取文件内容
		fileContent, err := ioutil.ReadFile(responseRule.File)
		if err != nil {
			log.Printf("读取文件失败: %v", err)
			// 如果文件读取失败，返回Body内容作为备选
			w.Write([]byte(responseRule.Body))
			log.Printf("返回响应: 状态码=%d, 内容类型=%s (文件读取失败，使用Body)", responseRule.Status, w.Header().Get("Content-Type"))
			log.Printf("响应内容: %s", responseRule.Body)
		} else {
			// 如果Content-Type未设置，根据文件扩展名自动设置
			if w.Header().Get("Content-Type") == "" {
				ext := filepath.Ext(responseRule.File)
				mimeType := mime.TypeByExtension(ext)
				if mimeType != "" {
					w.Header().Set("Content-Type", mimeType)
				} else {
					// 默认为二进制流
					w.Header().Set("Content-Type", "application/octet-stream")
				}
			}

			// 写入文件内容
			w.Write(fileContent)
			log.Printf("返回文件响应: 状态码=%d, 内容类型=%s, 文件=%s, 大小=%d字节",
				responseRule.Status, w.Header().Get("Content-Type"), responseRule.File, len(fileContent))
		}
	} else {
		// 写入响应体
		w.Write([]byte(responseRule.Body))
		log.Printf("返回响应: 状态码=%d, 内容类型=%s", responseRule.Status, w.Header().Get("Content-Type"))
		log.Printf("响应内容: %s", responseRule.Body)
	}

	log.Println("--------------------------------------------------")
}
