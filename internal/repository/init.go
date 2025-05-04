package repository

import (
	"ai-course/internal/config"
	"ai-course/internal/logger"

	"go.uber.org/zap"
)

// Initialize 初始化所有存储层组件
func Initialize(conf *config.MySQLConfig) error {
	// 检查是否配置了数据库凭据
	if conf.Username == "" || conf.Password == "" {
		logger.Logger.Warn("Skipping MySQL initialization - no credentials provided")
		return nil
	}

	// 初始化MySQL连接
	if err := InitMySQL(conf); err != nil {
		logger.Logger.Error("Failed to initialize MySQL", zap.Error(err))
		return err
	}

	// 执行数据库迁移
	if err := AutoMigrate(); err != nil {
		logger.Logger.Error("Failed to perform database migration", zap.Error(err))
		return err
	}

	logger.Logger.Info("Repository layer initialized successfully")
	return nil
}

// Cleanup 清理所有存储层资源
func Cleanup() error {
	if err := CloseMySQL(); err != nil {
		logger.Logger.Error("Failed to close MySQL connection", zap.Error(err))
		return err
	}

	logger.Logger.Info("Repository layer cleaned up successfully")
	return nil
}
