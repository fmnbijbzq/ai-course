package app

import (
	"ai-course/internal/config"
	"ai-course/internal/controller"
	"ai-course/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Application 应用程序结构
type Application struct {
	Engine       *gin.Engine
	Config       *config.Config
	UserService  service.UserService
	ClassService service.ClassService
}

// NewApplication 创建应用程序实例
func NewApplication(
	engine *gin.Engine,
	cfg *config.Config,
	userService service.UserService,
	classService service.ClassService,
) *Application {
	return &Application{
		Engine:       engine,
		Config:       cfg,
		UserService:  userService,
		ClassService: classService,
	}
}

// NewGinEngine 创建 Gin 引擎
func NewGinEngine(cfg *config.Config) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 gin 引擎
	engine := gin.New()
	return engine
}

// RegisterRoutes 注册路由
func (app *Application) RegisterRoutes() {
	router := controller.NewRouter(app.Engine, app.UserService, app.ClassService)
	router.RegisterRoutes()
}

// Run 运行应用程序
func (app *Application) Run() error {
	// 注册路由
	app.RegisterRoutes()

	// 启动服务器
	addr := fmt.Sprintf(":%d", app.Config.Server.Port)
	return app.Engine.Run(addr)
}
