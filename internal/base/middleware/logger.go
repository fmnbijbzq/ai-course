package middleware

import (
	"ai-course/internal/logger"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseWriter 是一个自定义的响应写入器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法以捕获响应内容
func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// APILogger 中间件函数，用于记录API请求和响应的详细信息
func APILogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取请求信息
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// 如果有查询参数，将其附加到路径
		if raw != "" {
			path = path + "?" + raw
		}

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为读取后会清空
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装ResponseWriter以捕获响应内容
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 获取响应状态
		status := c.Writer.Status()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		}

		// 添加请求头信息
		if userAgent := c.Request.UserAgent(); userAgent != "" {
			fields = append(fields, zap.String("user_agent", userAgent))
		}

		// 添加请求体（如果存在且不是二进制数据）
		if len(requestBody) > 0 && isTextContent(c.ContentType()) {
			fields = append(fields, zap.String("request_body", string(requestBody)))
		}

		// 添加响应体（如果存在且不是二进制数据）
		if w.body.Len() > 0 && isTextContent(c.Writer.Header().Get("Content-Type")) {
			fields = append(fields, zap.String("response_body", w.body.String()))
		}

		// 根据状态码选择日志级别
		switch {
		case status >= 500:
			logger.Logger.Error("API Request", fields...)
		case status >= 400:
			logger.Logger.Warn("API Request", fields...)
		default:
			logger.Logger.Info("API Request", fields...)
		}
	}
}

// isTextContent 检查内容类型是否为文本
func isTextContent(contentType string) bool {
	switch contentType {
	case "application/json",
		"application/xml",
		"application/x-www-form-urlencoded",
		"text/plain",
		"text/html",
		"text/xml",
		"text/json":
		return true
	default:
		return false
	}
}
