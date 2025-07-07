package controller

import (
	"ai-course/internal/base/controller"
	"ai-course/internal/logger"
	"ai-course/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserController 用户控制器
type UserController struct {
	controller.BaseController
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterRoutes 注册路由
func (c *UserController) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", c.Register)
		userGroup.POST("/login", c.Login)
	}
}

// Register godoc
// @Summary 用户注册
// @Description 注册新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.CreateUserDTO true "注册信息"
// @Success 200 {object} response.Response "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/user/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.CreateUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid register request",
			zap.Error(err),
		)
		c.ParamError("注册参数无效")
		return
	}

	err := c.userService.Register(ctx, &req)
	if err != nil {
		logger.Logger.Error("Failed to register user",
			zap.Error(err),
			zap.String("student_id", req.StudentID),
		)
		c.ServerError(err.Error())
		return
	}

	logger.Logger.Info("User registered successfully",
		zap.String("student_id", req.StudentID),
	)

	c.SuccessWithMessage("注册成功", nil)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录系统
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.LoginUserDTO true "登录信息"
// @Success 200 {object} response.Response "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Router /api/user/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req service.LoginUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid login request",
			zap.Error(err),
		)
		c.ParamError("登录参数无效")
		return
	}

	resp, err := c.userService.Login(ctx, &req)
	if err != nil {
		logger.Logger.Error("Failed to login",
			zap.Error(err),
			zap.String("student_id", req.StudentID),
		)
		c.Unauthorized(err.Error())
		return
	}

	logger.Logger.Info("User logged in successfully",
		zap.String("student_id", req.StudentID),
	)

	c.Success(resp)
}
