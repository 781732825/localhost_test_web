package utils

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dir string) error {
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// GetExecutableDir 获取可执行文件所在目录
func GetExecutableDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(execPath), nil
}

// WaitForSignal 等待系统信号（如Ctrl+C）
func WaitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}