package controller

import (
	"ai-course/internal/base/controller"
	"ai-course/internal/base/middleware"
	"ai-course/internal/logger"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	// 注册全局中间件
	r.Use(cors.New(corsConfig))
	r.Use(logger.GinZapLogger(), logger.GinZapRecovery(true))
	r.Use(middleware.APILogger()) // 添加API日志中间件

	// 创建基础控制器
	baseCtrl := &controller.BaseController{}

	// HealthCheck godoc
	// @Summary 健康检查
	// @Description 检查服务是否正常运行
	// @Tags 系统状态
	// @Produce json
	// @Success 200 {object} response.Response "服务正常运行"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		baseCtrl.InitHandler(c)
		baseCtrl.Success(gin.H{
			"status": "ok",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// ErrorTest godoc
	// @Summary 错误测试
	// @Description 测试错误处理
	// @Tags 系统状态
	// @Produce json
	// @Success 500 {object} response.Response "测试错误"
	// @Router /error [get]
	r.GET("/error", func(c *gin.Context) {
		baseCtrl.InitHandler(c)
		baseCtrl.ServerError("This is a test error")
	})

	// 用户路由组
	userController := NewUserController()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
	}

	// 班级路由组
	classController := NewClassController()
	classGroup := r.Group("/class")
	{
		classGroup.POST("/add", classController.Add)
		classGroup.PUT("/:id", classController.Edit)
		classGroup.DELETE("/:id", classController.Delete)
		classGroup.GET("/list", classController.List)
	}
}
