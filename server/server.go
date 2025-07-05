package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/trae/host/model"
)

// Server 表示HTTP服务器
type Server struct {
	config  *model.Config
	handler http.Handler
	servers []*http.Server
	wg      sync.WaitGroup
}

// NewServer 创建新的服务器
func NewServer(config *model.Config, handler http.Handler) *Server {
	return &Server{
		config:  config,
		handler: handler,
		servers: make([]*http.Server, 0),
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 遍历所有端口配置
	for _, portConfig := range s.config.Ports {
		// 创建HTTP服务器
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", portConfig.Port),
			Handler: s.handler,
		}

		// 保存服务器实例
		s.servers = append(s.servers, server)

		// 启动服务器
		s.wg.Add(1)
		go func(server *http.Server, portConfig model.PortConfig) {
			defer s.wg.Done()

			var err error
			if portConfig.HTTPS {
				// 启动HTTPS服务器
				log.Printf("启动HTTPS服务器，监听端口 %d", portConfig.Port)
				// 使用证书和密钥文件
				err = server.ListenAndServeTLS(portConfig.Cert, portConfig.Key)
			} else {
				// 启动HTTP服务器
				log.Printf("启动HTTP服务器，监听端口 %d", portConfig.Port)
				err = server.ListenAndServe()
			}

			if err != nil && err != http.ErrServerClosed {
				log.Printf("服务器错误: %v", err)
			}
		}(server, portConfig)
	}

	return nil
}



// Stop 停止服务器
func (s *Server) Stop() {
	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 停止所有服务器
	for _, server := range s.servers {
		server.Shutdown(ctx)
	}

	// 等待所有服务器停止
	s.wg.Wait()
	log.Println("所有服务器已停止")
}