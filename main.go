package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/trae/host/config"
	"github.com/trae/host/handler"
	"github.com/trae/host/server"
	"github.com/trae/host/utils"
)

func main() {
	// 获取可执行文件所在目录
	execDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatalf("获取可执行文件目录失败: %v", err)
	}

	// 设置默认配置文件和规则目录路径（相对于可执行文件目录）
	defaultConfigPath := filepath.Join(execDir, "config.json")
	defaultRuleDir := filepath.Join(execDir, "rules")

	// 解析命令行参数
	configPath := flag.String("config", "", "配置文件路径")
	ruleDir := flag.String("rules", "", "规则文件目录")
	flag.Parse()

	// 检查是否指定了配置文件路径
	if *configPath == "" {
		log.Printf("未指定配置文件路径，将使用默认路径: %s", defaultConfigPath)
		*configPath = defaultConfigPath
	} else if !filepath.IsAbs(*configPath) {
		// 如果配置文件路径是相对路径，则相对于可执行文件目录
		*configPath = filepath.Join(execDir, *configPath)
	}

	// 检查是否指定了规则目录
	if *ruleDir == "" {
		log.Printf("未指定规则目录，将使用默认目录: %s", defaultRuleDir)
		*ruleDir = defaultRuleDir
	} else if !filepath.IsAbs(*ruleDir) {
		// 如果规则目录是相对路径，则相对于可执行文件目录
		*ruleDir = filepath.Join(execDir, *ruleDir)
	}

	// 确保规则目录存在
	if err := utils.EnsureDir(*ruleDir); err != nil {
		log.Fatalf("创建规则目录失败: %v", err)
	}

	// 加载配置
	log.Printf("加载配置文件: %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建规则管理器
	ruleManager := handler.NewRuleManager(*ruleDir, cfg.DefaultResponse)

	// 创建请求处理器
	requestHandler := handler.NewRequestHandler(ruleManager)

	// 创建服务器
	srv := server.NewServer(cfg, requestHandler)

	// 启动服务器
	log.Println("#######################################")
	log.Println("本地调试Web服务器")
	log.Println("By: 简单")
	log.Println("项目地址: https://github.com/781732825/localhost_test_web")
	log.Println("警告：本工具仅限于 本地 调试/测试 使用，请按照自行编写相关规则，用于验证数据是否正常。")
	log.Println("警告：请勿在 公网 环境使用，否则后果自负。")
	log.Println("警告：请勿将 工具 用于任何非法用途，否则后果自负。")
	log.Println("#######################################")

	log.Println("正在启动服务器...")
	if err := srv.Start(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	// 等待信号
	log.Println("服务器已启动，按Ctrl+C停止")
	utils.WaitForSignal()

	// 停止服务器
	log.Println("正在停止服务器...")
	srv.Stop()

	log.Println("服务器已停止")
	os.Exit(0)
}
