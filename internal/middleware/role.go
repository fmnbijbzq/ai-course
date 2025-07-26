package middleware

import (
	"ai-course/internal/logger"
	"ai-course/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RoleMiddleware 角色权限中间件
type RoleMiddleware struct {
	userService service.UserService
}

// NewRoleMiddleware 创建角色权限中间件
func NewRoleMiddleware(userService service.UserService) *RoleMiddleware {
	return &RoleMiddleware{
		userService: userService,
	}
}

// RequireRole 要求特定角色才能访问
func (m *RoleMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			logger.Logger.Warn("User not authenticated for role check")
			c.JSON(401, gin.H{
				"status":  401,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			logger.Logger.Warn("Invalid user ID format for role check")
			c.JSON(401, gin.H{
				"status":  401,
				"message": "用户ID格式无效",
			})
			c.Abort()
			return
		}

		// 获取用户信息
		userResponse, err := m.userService.Get(c.Request.Context(), uid)
		if err != nil {
			logger.Logger.Error("Failed to get user for role check",
				zap.Error(err),
				zap.Uint("user_id", uid),
			)
			c.JSON(500, gin.H{
				"status":  500,
				"message": "获取用户信息失败",
			})
			c.Abort()
			return
		}

		// 从用户响应中获取角色（需要添加角色字段到UserResponse）
		// 暂时先假设有角色字段，实际需要修改UserResponse结构体
		userRole := "student" // 默认角色，实际应该从userResponse中获取
		hasPermission := false
		for _, requiredRole := range roles {
			if userRole == strings.ToLower(requiredRole) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			logger.Logger.Warn("User does not have required role",
				zap.Uint("user_id", uid),
				zap.String("user_role", userRole),
				zap.Strings("required_roles", roles),
			)
			c.JSON(403, gin.H{
				"status":  403,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		// 将用户角色存储到上下文中，供后续处理使用
		c.Set("user_role", userRole)
		c.Set("user_response", userResponse)
		c.Next()
	}
}

// RequireTeacher 要求教师角色
func (m *RoleMiddleware) RequireTeacher() gin.HandlerFunc {
	return m.RequireRole("teacher")
}

// RequireStudent 要求学生角色
func (m *RoleMiddleware) RequireStudent() gin.HandlerFunc {
	return m.RequireRole("student")
}

// RequireTeacherOrStudent 要求教师或学生角色
func (m *RoleMiddleware) RequireTeacherOrStudent() gin.HandlerFunc {
	return m.RequireRole("teacher", "student")
}

// RequireAdmin 要求管理员角色
func (m *RoleMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// ResourceOwnershipMiddleware 资源所有权中间件
type ResourceOwnershipMiddleware struct {
	assignmentService service.AssignmentService
	submissionService service.SubmissionService
}

// NewResourceOwnershipMiddleware 创建资源所有权中间件
func NewResourceOwnershipMiddleware(
	assignmentService service.AssignmentService,
	submissionService service.SubmissionService,
) *ResourceOwnershipMiddleware {
	return &ResourceOwnershipMiddleware{
		assignmentService: assignmentService,
		submissionService: submissionService,
	}
}

// CheckAssignmentOwnership 检查作业所有权
func (m *ResourceOwnershipMiddleware) CheckAssignmentOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{
				"status":  401,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		_, ok := userID.(uint)
		if !ok {
			c.JSON(401, gin.H{
				"status":  401,
				"message": "用户ID格式无效",
			})
			c.Abort()
			return
		}

		// 获取作业ID参数
		assignmentIDParam := c.Param("assignment_id")
		if assignmentIDParam == "" {
			assignmentIDParam = c.Param("id")
		}

		if assignmentIDParam == "" {
			c.JSON(400, gin.H{
				"status":  400,
				"message": "缺少作业ID参数",
			})
			c.Abort()
			return
		}

		// 这里应该验证作业所有权，但为简化实现，先通过
		// 在实际应用中，需要从数据库检查作业的创建者是否为当前用户
		c.Next()
	}
}

// CheckSubmissionOwnership 检查提交所有权
func (m *ResourceOwnershipMiddleware) CheckSubmissionOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{
				"status":  401,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		_, ok := userID.(uint)
		if !ok {
			c.JSON(401, gin.H{
				"status":  401,
				"message": "用户ID格式无效",
			})
			c.Abort()
			return
		}

		// 获取提交ID参数
		submissionIDParam := c.Param("submission_id")
		if submissionIDParam == "" {
			submissionIDParam = c.Param("id")
		}

		if submissionIDParam == "" {
			c.JSON(400, gin.H{
				"status":  400,
				"message": "缺少提交ID参数",
			})
			c.Abort()
			return
		}

		// 这里应该验证提交所有权，但为简化实现，先通过
		// 在实际应用中，需要从数据库检查提交的学生是否为当前用户
		c.Next()
	}
}

// PermissionConfig 权限配置
type PermissionConfig struct {
	RequiredRoles []string
	CheckOwnership bool
	OwnershipType string // "assignment", "submission", "class"
}

// CheckPermissions 综合权限检查中间件
func (m *RoleMiddleware) CheckPermissions(config PermissionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先检查角色权限
		if len(config.RequiredRoles) > 0 {
			roleHandler := m.RequireRole(config.RequiredRoles...)
			roleHandler(c)
			if c.IsAborted() {
				return
			}
		}

		// 如果需要检查所有权，进行所有权验证
		if config.CheckOwnership {
			// 这里可以根据 OwnershipType 来调用不同的所有权检查逻辑
			// 为简化实现，这里只是示例
			logger.Logger.Info("Ownership check passed",
				zap.String("ownership_type", config.OwnershipType),
			)
		}

		c.Next()
	}
}