package main

import (
	"ai-course/docs" // 导入 swagger 文档
	"ai-course/internal/config"
	"ai-course/internal/handler"
	"ai-course/internal/logger"
	"ai-course/internal/repository"
	"ai-course/internal/server"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AI Course API
// @version 1.0
// @description This is the API documentation for AI Course Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// 初始化 swagger 文档
	docs.SwaggerInfo.BasePath = "/"

	// 初始化配置
	config.InitConfig()

	// 创建日志目录
	logDir := filepath.Dir(config.GlobalConfig.Logger.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("create log directory failed: %v", err))
	}

	// 初始化日志
	logger.InitLogger(&config.GlobalConfig.Logger)
	defer logger.Logger.Sync()

	// 初始化存储层
	if err := repository.Initialize(&config.GlobalConfig.MySQL); err != nil {
		logger.Logger.Fatal("Failed to initialize repository layer")
	}
	defer repository.Cleanup()

	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建gin引擎
	engine := gin.New()

	// 注册路由
	handler.RegisterRoutes(engine)

	// 添加Swagger路由
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建并启动服务器
	srv := server.NewServer(engine, &config.GlobalConfig.Server)
	if err := srv.Start(); err != nil {
		logger.Logger.Fatal("Server failed to start")
	}
}
