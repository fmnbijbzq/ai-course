package repository

import (
	"ai-course/internal/logger"
	"ai-course/internal/model"

	"go.uber.org/zap"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	logger.Logger.Info("Starting database migration...")

	// 在这里添加需要迁移的模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.Class{},
		// 添加其他模型...
	)

	if err != nil {
		logger.Logger.Error("Database migration failed", zap.Error(err))
		return err
	}

	logger.Logger.Info("Database migration completed successfully")
	return nil
}
