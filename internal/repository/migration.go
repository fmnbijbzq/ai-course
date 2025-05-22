package repository

import (
	"ai-course/internal/logger"
	"ai-course/internal/model"
	"context"

	"go.uber.org/zap"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db DB) error {
	logger.Logger.Info("Starting database migration...")

	err := db.WithContext(context.Background()).AutoMigrate(
		&model.User{},
		&model.Class{},
	)

	if err != nil {
		logger.Logger.Error("Database migration failed", zap.Error(err))
		return err
	}

	logger.Logger.Info("Database migration completed successfully")
	return nil
}
