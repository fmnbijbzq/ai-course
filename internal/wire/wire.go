//go:build wireinject
// +build wireinject

package wire

import (
	"ai-course/internal/app"
	"ai-course/internal/config"
	"ai-course/internal/repository"
	"ai-course/internal/service"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// InitializeApplication 初始化应用程序
func InitializeApplication() (*app.Application, error) {
	wire.Build(
		// 配置初始化
		config.LoadConfig,

		// 数据库初始化
		config.InitDB,

		// Repository 层
		repository.NewGormDB,
		repository.NewNoOpCache,
		repository.NewUserRepository,
		repository.NewClassRepository,

		// Service 层
		service.NewUserService,
		service.NewClassService,

		// Gin 引擎
		app.NewGinEngine,

		// 应用程序
		app.NewApplication,
	)
	return nil, nil
}

// InitializeUserService 初始化用户服务（保留用于兼容性）
func InitializeUserService(db *gorm.DB) (service.UserService, error) {
	wire.Build(
		repository.NewGormDB,
		repository.NewNoOpCache,
		repository.NewUserRepository,
		service.NewUserService,
	)
	return nil, nil
}

// InitializeClassService 初始化班级服务
func InitializeClassService(db *gorm.DB) (service.ClassService, error) {
	wire.Build(
		repository.NewGormDB,
		repository.NewNoOpCache,
		repository.NewClassRepository,
		service.NewClassService,
	)
	return nil, nil
}
