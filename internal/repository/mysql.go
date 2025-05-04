package repository

import (
	"ai-course/internal/config"
	"ai-course/internal/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormZapLogger GORM的Zap日志适配器
type GormZapLogger struct {
	SlowThreshold        time.Duration
	IgnoreRecordNotFound bool
}

func NewGormZapLogger() *GormZapLogger {
	return &GormZapLogger{
		SlowThreshold:        time.Second, // 慢SQL阈值
		IgnoreRecordNotFound: true,
	}
}

// LogMode 实现 gormlogger.Interface 接口
func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info 实现 gormlogger.Interface 接口
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.Logger.Info(fmt.Sprintf(msg, data...),
		zap.String("type", "gorm"),
	)
}

// Warn 实现 gormlogger.Interface 接口
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.Logger.Warn(fmt.Sprintf(msg, data...),
		zap.String("type", "gorm"),
	)
}

// Error 实现 gormlogger.Interface 接口
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Logger.Error(fmt.Sprintf(msg, data...),
		zap.String("type", "gorm"),
	)
}

// Trace 实现 gormlogger.Interface 接口
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("type", "gorm"),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFound) {
		fields = append(fields, zap.Error(err))
		logger.Logger.Error("SQL Error", fields...)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger.Logger.Warn("Slow SQL", fields...)
		return
	}

	logger.Logger.Debug("SQL Trace", fields...)
}

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL(conf *config.MySQLConfig) error {
	gormLogger := NewGormZapLogger()

	db, err := gorm.Open(mysql.Open(conf.GetMySQLDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)

	DB = db
	return nil
}

// CloseMySQL 关闭MySQL连接
func CloseMySQL() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
