package main

import (
	"ai-course/docs"
	"ai-course/internal/logger"
	"ai-course/internal/wire"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AI Course API
// @version 1.0
// @description AI Course 后端 API 服务
// @BasePath /
func main() {
	// 初始化 Swagger 文档
	docs.SwaggerInfo.Title = "AI Course API"
	docs.SwaggerInfo.Description = "AI Course 后端 API 服务"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 使用 wire 初始化应用程序
	app, err := wire.InitializeApplication()
	if err != nil {
		panic("Failed to initialize application: " + err.Error())
	}

	// 初始化日志
	logger.InitLogger(&app.Config.Logger)

	// 添加 Swagger 路由
	app.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 运行应用程序
	if err := app.Run(); err != nil {
		logger.Logger.Fatal("Failed to start server: " + err.Error())
	}
}
