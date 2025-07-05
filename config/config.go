package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/trae/host/model"
)

// LoadConfig 从指定路径加载配置文件
func LoadConfig(path string) (*model.Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 如果配置文件不存在，创建默认配置
		return createDefaultConfig(path)
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON
	var config model.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(path string) (*model.Config, error) {
	// 创建默认配置
	config := &model.Config{
		Ports: []model.PortConfig{
			{
				Port:  8080,
				HTTPS: false,
			},
		},
		DefaultResponse: "没有匹配的规则",
	}

	// 将配置写入文件
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return nil, fmt.Errorf("写入配置文件失败: %v", err)
	}

	return config, nil
}

// validateConfig 验证配置是否有效
func validateConfig(config *model.Config) error {
	// 检查是否有端口配置
	if len(config.Ports) == 0 {
		return fmt.Errorf("配置文件中没有端口配置")
	}

	// 检查HTTPS配置
	for _, port := range config.Ports {
		if port.HTTPS {
			// 检查证书和密钥文件
			if port.Cert == "" || port.Key == "" {
				return fmt.Errorf("HTTPS端口 %d 缺少证书或密钥配置", port.Port)
			}

			// 检查证书文件是否存在
			if _, err := os.Stat(port.Cert); os.IsNotExist(err) {
				return fmt.Errorf("证书文件不存在: %s", port.Cert)
			}

			// 检查密钥文件是否存在
			if _, err := os.Stat(port.Key); os.IsNotExist(err) {
				return fmt.Errorf("密钥文件不存在: %s", port.Key)
			}
		}
	}

	return nil
}