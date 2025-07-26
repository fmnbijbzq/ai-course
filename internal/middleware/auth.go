package middleware

import (
	"ai-course/internal/logger"
	"ai-course/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Logger.Warn("Missing authorization header")
			c.JSON(401, gin.H{
				"status":  401,
				"message": "缺少认证信息",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Logger.Warn("Invalid authorization header format")
			c.JSON(401, gin.H{
				"status":  401,
				"message": "认证信息格式无效",
			})
			c.Abort()
			return
		}

		// 提取token
		tokenString := authHeader[7:] // 去除"Bearer "前缀
		if tokenString == "" {
			logger.Logger.Warn("Empty token")
			c.JSON(401, gin.H{
				"status":  401,
				"message": "认证令牌为空",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			logger.Logger.Warn("Invalid token",
				zap.Error(err),
				zap.String("token", tokenString[:min(len(tokenString), 20)]+"..."),
			)
			c.JSON(401, gin.H{
				"status":  401,
				"message": "认证令牌无效",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("student_id", claims.StudentID)

		logger.Logger.Debug("User authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("student_id", claims.StudentID),
		)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有认证信息，继续处理
			c.Next()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// 格式无效，继续处理
			c.Next()
			return
		}

		// 提取token
		tokenString := authHeader[7:] // 去除"Bearer "前缀
		if tokenString == "" {
			// token为空，继续处理
			c.Next()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			// token无效，继续处理
			logger.Logger.Debug("Optional auth failed, continuing without auth",
				zap.Error(err),
			)
			c.Next()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("student_id", claims.StudentID)

		logger.Logger.Debug("User optionally authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("student_id", claims.StudentID),
		)

		c.Next()
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}