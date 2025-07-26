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

	// 在迁移前清理可能的重复空值数据
	err := cleanupDuplicateData(db)
	if err != nil {
		logger.Logger.Error("Failed to cleanup duplicate data", zap.Error(err))
		return err
	}

	err = db.WithContext(context.Background()).AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Class{},
		&model.Assignment{},
		&model.Question{},
		&model.Submission{},
		&model.Answer{},
		&model.Attachment{},
	)

	if err != nil {
		logger.Logger.Error("Database migration failed", zap.Error(err))
		return err
	}

	logger.Logger.Info("Database migration completed successfully")
	return nil
}

// cleanupDuplicateData 清理重复的空值数据
func cleanupDuplicateData(db DB) error {
	logger.Logger.Info("Cleaning up duplicate data...")

	// 检查users表是否存在
	var tableExists int
	err := db.WithContext(context.Background()).Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() 
		AND table_name = 'users'
	`).Scan(&tableExists)
	
	if err != nil {
		logger.Logger.Error("Failed to check if users table exists", zap.Error(err))
		return err
	}

	if tableExists == 0 {
		logger.Logger.Info("Users table does not exist, skipping cleanup")
		return nil
	}

	// 检查是否存在空code的记录
	var emptyCodeCount int
	err = db.WithContext(context.Background()).Raw(`
		SELECT COUNT(*) 
		FROM users 
		WHERE code = '' OR code IS NULL
	`).Scan(&emptyCodeCount)
	
	if err != nil {
		logger.Logger.Error("Failed to check empty code records", zap.Error(err))
		return err
	}

	if emptyCodeCount == 0 {
		logger.Logger.Info("No empty code records found, skipping cleanup")
		return nil
	}

	logger.Logger.Info("Found empty code records, cleaning up...", zap.Int("count", emptyCodeCount))

	// 为空code记录设置唯一值
	err = db.WithContext(context.Background()).Exec(`
		UPDATE users 
		SET code = CONCAT('temp_user_', id, '_', UNIX_TIMESTAMP(NOW(3)))
		WHERE code = '' OR code IS NULL
	`)

	if err != nil {
		logger.Logger.Error("Failed to update empty code records", zap.Error(err))
		return err
	}

	logger.Logger.Info("Duplicate data cleanup completed successfully")
	return nil
}
