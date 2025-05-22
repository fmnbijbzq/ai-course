package controller

import (
	"ai-course/internal/base/response"

	"github.com/gin-gonic/gin"
)

// BaseController 基础控制器
type BaseController struct {
	handler *response.Handler
}

// InitHandler 初始化响应处理器
func (b *BaseController) InitHandler(c *gin.Context) {
	b.handler = response.NewHandler(c)
}

// Success 成功响应
func (b *BaseController) Success(data interface{}) {
	b.handler.Success(data)
}

// SuccessWithMessage 自定义消息的成功响应
func (b *BaseController) SuccessWithMessage(message string, data interface{}) {
	b.handler.SuccessWithMessage(message, data)
}

// Fail 失败响应
func (b *BaseController) Fail(code int, message string) {
	b.handler.Fail(code, message)
}

// ParamError 参数错误响应
func (b *BaseController) ParamError(message string) {
	b.handler.ParamError(message)
}

// Unauthorized 未授权响应
func (b *BaseController) Unauthorized(message string) {
	b.handler.Unauthorized(message)
}

// ServerError 服务器错误响应
func (b *BaseController) ServerError(message string) {
	b.handler.ServerError(message)
}
