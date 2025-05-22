package main

import (
	"ai-course/internal/logger"
	"ai-course/internal/wire"
)

func main() {
	// 使用 wire 初始化应用程序
	app, err := wire.InitializeApplication()
	if err != nil {
		panic("Failed to initialize application: " + err.Error())
	}

	// 初始化日志
	logger.InitLogger(&app.Config.Logger)

	// 运行应用程序
	if err := app.Run(); err != nil {
		logger.Logger.Fatal("Failed to start server: " + err.Error())
	}
}
