// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"ai-course/internal/app"
	"ai-course/internal/config"
	"ai-course/internal/repository"
	"ai-course/internal/service"
	"gorm.io/gorm"
)

// Injectors from wire.go:

// InitializeApplication 初始化应用程序
func InitializeApplication() (*app.Application, error) {
	configConfig := config.LoadConfig()
	engine := app.NewGinEngine(configConfig)
	db := config.InitDB(configConfig)
	repositoryDB := repository.NewGormDB(db)
	cache := repository.NewNoOpCache()
	userRepository := repository.NewUserRepository(repositoryDB, cache)
	userService := service.NewUserService(userRepository)
	classRepository := repository.NewClassRepository(repositoryDB, cache)
	classService := service.NewClassService(classRepository)
	application := app.NewApplication(engine, configConfig, userService, classService)
	return application, nil
}

// InitializeUserService 初始化用户服务（保留用于兼容性）
func InitializeUserService(db *gorm.DB) (service.UserService, error) {
	repositoryDB := repository.NewGormDB(db)
	cache := repository.NewNoOpCache()
	userRepository := repository.NewUserRepository(repositoryDB, cache)
	userService := service.NewUserService(userRepository)
	return userService, nil
}

// InitializeClassService 初始化班级服务
func InitializeClassService(db *gorm.DB) (service.ClassService, error) {
	repositoryDB := repository.NewGormDB(db)
	cache := repository.NewNoOpCache()
	classRepository := repository.NewClassRepository(repositoryDB, cache)
	classService := service.NewClassService(classRepository)
	return classService, nil
}
