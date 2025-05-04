package handler

import (
	"ai-course/internal/logger"
	"ai-course/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Register godoc
// @Summary 用户注册
// @Description 注册新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.UserRegisterRequest true "注册信息"
// @Success 200 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid register request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	resp, err := h.userService.Register(&req)
	if err != nil {
		logger.Logger.Error("Failed to register user",
			zap.Error(err),
			zap.String("student_id", req.StudentID),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Logger.Info("User registered successfully",
		zap.String("student_id", req.StudentID),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user":    resp,
	})
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录系统
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.UserLoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Router /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid login request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		logger.Logger.Error("Failed to login",
			zap.Error(err),
			zap.String("student_id", req.StudentID),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.Logger.Info("User logged in successfully",
		zap.String("student_id", req.StudentID),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    resp,
	})
}
