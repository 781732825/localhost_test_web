package model

// PortConfig 表示单个端口的配置信息
type PortConfig struct {
	Port        int    `json:"port"`         // 监听的端口号
	HTTPS       bool   `json:"https"`        // 是否启用HTTPS
	Cert        string `json:"cert"`         // 证书文件路径（仅HTTPS时有效）
	Key         string `json:"key"`          // 密钥文件路径（仅HTTPS时有效）
}

// Config 表示整个应用的配置信息
type Config struct {
	Ports           []PortConfig `json:"ports"`           // 端口配置列表
	DefaultResponse string       `json:"defaultResponse"` // 默认响应内容
}