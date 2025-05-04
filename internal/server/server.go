package server

import (
	"ai-course/internal/config"
	"ai-course/internal/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server
}

// NewServer 创建新的服务器实例
func NewServer(engine *gin.Engine, conf *config.ServerConfig) *Server {
	return &Server{
		engine: engine,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.Port),
			Handler: engine,
		},
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 在goroutine中启动服务器
	go func() {
		logger.Logger.Info("Starting server...",
			zap.String("addr", s.srv.Addr),
		)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down server...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	logger.Logger.Info("Server exiting")
	return nil
}
