package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 响应信息
	Data    interface{} `json:"data"`    // 响应数据
}

// ResponseCode 定义业务状态码
const (
	CodeSuccess      = 200 // 成功
	CodeInvalidParam = 400 // 请求参数错误
	CodeUnauthorized = 401 // 未授权
	CodeForbidden    = 403 // 禁止访问
	CodeNotFound     = 404 // 资源不存在
	CodeServerError  = 500 // 服务器内部错误
)

// 预定义错误消息
const (
	MsgSuccess      = "操作成功"
	MsgInvalidParam = "请求参数错误"
	MsgUnauthorized = "未授权访问"
	MsgForbidden    = "禁止访问"
	MsgNotFound     = "资源不存在"
	MsgServerError  = "服务器内部错误"
)

// Handler 响应处理器
type Handler struct {
	C *gin.Context
}

// NewHandler 创建响应处理器
func NewHandler(c *gin.Context) *Handler {
	return &Handler{C: c}
}

// Success 成功响应
func (r *Handler) Success(data interface{}) {
	resp := Response{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
	}
	r.C.JSON(http.StatusOK, resp)
}

// SuccessWithMessage 自定义消息的成功响应
func (r *Handler) SuccessWithMessage(message string, data interface{}) {
	resp := Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
	r.C.JSON(http.StatusOK, resp)
}

// Fail 失败响应
func (r *Handler) Fail(code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	r.C.JSON(http.StatusOK, resp)
}

// FailWithData 带数据的失败响应
func (r *Handler) FailWithData(code int, message string, data interface{}) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	r.C.JSON(http.StatusOK, resp)
}

// ParamError 参数错误响应
func (r *Handler) ParamError(message string) {
	if message == "" {
		message = MsgInvalidParam
	}
	resp := Response{
		Code:    CodeInvalidParam,
		Message: message,
		Data:    nil,
	}
	r.C.JSON(http.StatusBadRequest, resp)
}

// Unauthorized 未授权响应
func (r *Handler) Unauthorized(message string) {
	if message == "" {
		message = MsgUnauthorized
	}
	resp := Response{
		Code:    CodeUnauthorized,
		Message: message,
		Data:    nil,
	}
	r.C.JSON(http.StatusUnauthorized, resp)
}

// Forbidden 禁止访问响应
func (r *Handler) Forbidden(message string) {
	if message == "" {
		message = MsgForbidden
	}
	resp := Response{
		Code:    CodeForbidden,
		Message: message,
		Data:    nil,
	}
	r.C.JSON(http.StatusForbidden, resp)
}

// NotFound 资源不存在响应
func (r *Handler) NotFound(message string) {
	if message == "" {
		message = MsgNotFound
	}
	resp := Response{
		Code:    CodeNotFound,
		Message: message,
		Data:    nil,
	}
	r.C.JSON(http.StatusNotFound, resp)
}

// ServerError 服务器错误响应
func (r *Handler) ServerError(message string) {
	if message == "" {
		message = MsgServerError
	}
	resp := Response{
		Code:    CodeServerError,
		Message: message,
		Data:    nil,
	}
	r.C.JSON(http.StatusInternalServerError, resp)
}

// JSON 发送JSON响应
func (r *Handler) JSON(code int, message string, data interface{}) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	r.C.JSON(http.StatusOK, resp)
}

// CustomJSON 自定义HTTP状态码的JSON响应
func (r *Handler) CustomJSON(httpStatus int, code int, message string, data interface{}) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	r.C.JSON(httpStatus, resp)
}
