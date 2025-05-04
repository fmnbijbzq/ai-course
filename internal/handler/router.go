package handler

import (
	"ai-course/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 注册中间件
	r.Use(logger.GinZapLogger(), logger.GinZapRecovery(true))

	// HealthCheck godoc
	// @Summary 健康检查
	// @Description 检查服务是否正常运行
	// @Tags 系统状态
	// @Produce json
	// @Success 200 {object} map[string]interface{} "服务正常运行"
	// @Router /health [get]
	r.GET("/health", HealthCheck)

	// ErrorTest godoc
	// @Summary 错误测试
	// @Description 测试错误处理
	// @Tags 系统状态
	// @Produce json
	// @Success 500 {object} map[string]interface{} "测试错误"
	// @Router /error [get]
	r.GET("/error", ErrorTest)

	// 用户路由组
	userHandler := NewUserHandler()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}

	// 班级路由组
	classHandler := NewClassHandler()
	classGroup := r.Group("/class")
	{
		classGroup.POST("/add", classHandler.Add)
		classGroup.PUT("/:id", classHandler.Edit)
		classGroup.DELETE("/:id", classHandler.Delete)
		classGroup.GET("/list", classHandler.List)
	}
}

// HealthCheck 健康检查处理函数
func HealthCheck(c *gin.Context) {
	logger.Logger.Info("Health check endpoint called",
		zap.String("remote_addr", c.Request.RemoteAddr),
		zap.String("user_agent", c.Request.UserAgent()),
	)
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   time.Now().Format("2006-01-02 15:04:05"),
	})
}

// ErrorTest 错误测试处理函数
func ErrorTest(c *gin.Context) {
	logger.Logger.Error("This is a test error",
		zap.String("custom_field", "test value"),
		zap.Int("status_code", 500),
	)
	c.JSON(500, gin.H{
		"status":  "error",
		"message": "This is a test error",
	})
}
