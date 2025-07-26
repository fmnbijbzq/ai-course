package controller

import (
	"ai-course/internal/base/controller"
	"ai-course/internal/logger"
	"ai-course/internal/model"
	"ai-course/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AssignmentController 作业控制器
type AssignmentController struct {
	controller.BaseController
	assignmentService service.AssignmentService
}

// NewAssignmentController 创建作业控制器
func NewAssignmentController(assignmentService service.AssignmentService) *AssignmentController {
	return &AssignmentController{
		assignmentService: assignmentService,
	}
}

// Create godoc
// @Summary 创建作业
// @Description 教师创建新作业
// @Tags 作业管理
// @Accept json
// @Produce json
// @Param request body model.CreateAssignmentRequest true "作业信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment [post]
func (c *AssignmentController) Create(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req model.CreateAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid create assignment request",
			zap.Error(err),
		)
		c.ParamError("创建作业参数无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for create assignment")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for create assignment")
		c.Unauthorized("用户ID格式无效")
		return
	}

	assignment, err := c.assignmentService.CreateAssignment(ctx.Request.Context(), &req, teacherID)
	if err != nil {
		logger.Logger.Error("Failed to create assignment",
			zap.Error(err),
			zap.Uint("teacher_id", teacherID),
		)
		c.ServerError(err.Error())
		return
	}

	logger.Logger.Info("Assignment created successfully",
		zap.Uint("assignment_id", assignment.ID),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("创建作业成功", assignment)
}

// Update godoc
// @Summary 更新作业
// @Description 教师更新作业信息
// @Tags 作业管理
// @Accept json
// @Produce json
// @Param id path int true "作业ID"
// @Param request body model.UpdateAssignmentRequest true "作业信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{id} [put]
func (c *AssignmentController) Update(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	var req model.UpdateAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid update assignment request",
			zap.Error(err),
		)
		c.ParamError("更新作业参数无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for update assignment")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for update assignment")
		c.Unauthorized("用户ID格式无效")
		return
	}

	assignment, err := c.assignmentService.UpdateAssignment(ctx.Request.Context(), uint(id), &req, teacherID)
	if err != nil {
		logger.Logger.Error("Failed to update assignment",
			zap.Error(err),
			zap.Uint("assignment_id", uint(id)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to update this assignment":
			c.Fail(403, "无权限操作此作业")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Assignment updated successfully",
		zap.Uint("assignment_id", assignment.ID),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("更新作业成功", assignment)
}

// Delete godoc
// @Summary 删除作业
// @Description 教师删除作业
// @Tags 作业管理
// @Produce json
// @Param id path int true "作业ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{id} [delete]
func (c *AssignmentController) Delete(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for delete assignment")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for delete assignment")
		c.Unauthorized("用户ID格式无效")
		return
	}

	err = c.assignmentService.DeleteAssignment(ctx.Request.Context(), uint(id), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to delete assignment",
			zap.Error(err),
			zap.Uint("assignment_id", uint(id)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to delete this assignment":
			c.Fail(403, "无权限操作此作业")
		case "cannot delete assignment with submissions":
			c.Fail(400, "已有学生提交的作业不能删除")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Assignment deleted successfully",
		zap.Uint("assignment_id", uint(id)),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("删除作业成功", nil)
}

// Detail godoc
// @Summary 获取作业详情
// @Description 获取作业详细信息
// @Tags 作业管理
// @Produce json
// @Param id path int true "作业ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{id} [get]
func (c *AssignmentController) Detail(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	assignment, err := c.assignmentService.GetAssignmentDetail(ctx.Request.Context(), uint(id))
	if err != nil {
		logger.Logger.Error("Failed to get assignment detail",
			zap.Error(err),
			zap.Uint("assignment_id", uint(id)),
		)
		
		if err.Error() == "assignment not found" {
			c.Fail(404, "作业不存在")
			return
		}
		c.ServerError(err.Error())
		return
	}

	c.Success(assignment)
}

// List godoc
// @Summary 获取教师作业列表
// @Description 获取当前教师的所有作业列表
// @Tags 作业管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/list [get]
func (c *AssignmentController) List(ctx *gin.Context) {
	c.InitHandler(ctx)
	
	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for list assignments")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for list assignments")
		c.Unauthorized("用户ID格式无效")
		return
	}

	assignments, total, err := c.assignmentService.GetTeacherAssignments(ctx.Request.Context(), teacherID, page, pageSize)
	if err != nil {
		logger.Logger.Error("Failed to get teacher assignments",
			zap.Error(err),
			zap.Uint("teacher_id", teacherID),
		)
		c.ServerError(err.Error())
		return
	}

	response := gin.H{
		"list":  assignments,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}

	c.Success(response)
}

// Publish godoc
// @Summary 发布作业
// @Description 教师发布作业
// @Tags 作业管理
// @Produce json
// @Param id path int true "作业ID"
// @Success 200 {object} response.Response "发布成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{id}/publish [post]
func (c *AssignmentController) Publish(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for publish assignment")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for publish assignment")
		c.Unauthorized("用户ID格式无效")
		return
	}

	err = c.assignmentService.PublishAssignment(ctx.Request.Context(), uint(id), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to publish assignment",
			zap.Error(err),
			zap.Uint("assignment_id", uint(id)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to publish this assignment":
			c.Fail(403, "无权限操作此作业")
		case "assignment is already published":
			c.Fail(400, "作业已经发布")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Assignment published successfully",
		zap.Uint("assignment_id", uint(id)),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("发布作业成功", nil)
}

// Unpublish godoc
// @Summary 取消发布作业
// @Description 教师取消发布作业
// @Tags 作业管理
// @Produce json
// @Param id path int true "作业ID"
// @Success 200 {object} response.Response "取消发布成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{id}/unpublish [post]
func (c *AssignmentController) Unpublish(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for unpublish assignment")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for unpublish assignment")
		c.Unauthorized("用户ID格式无效")
		return
	}

	err = c.assignmentService.UnpublishAssignment(ctx.Request.Context(), uint(id), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to unpublish assignment",
			zap.Error(err),
			zap.Uint("assignment_id", uint(id)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to unpublish this assignment":
			c.Fail(403, "无权限操作此作业")
		case "cannot unpublish assignment with submissions":
			c.Fail(400, "已有学生提交的作业不能取消发布")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Assignment unpublished successfully",
		zap.Uint("assignment_id", uint(id)),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("取消发布作业成功", nil)
}

// Statistics godoc
// @Summary 获取作业统计信息
// @Description 教师获取作业统计信息
// @Tags 作业管理
// @Produce json
// @Param id path int true "作业ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{id}/statistics [get]
func (c *AssignmentController) Statistics(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for get assignment statistics")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for get assignment statistics")
		c.Unauthorized("用户ID格式无效")
		return
	}

	stats, err := c.assignmentService.GetAssignmentStatistics(ctx.Request.Context(), uint(id), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to get assignment statistics",
			zap.Error(err),
			zap.Uint("assignment_id", uint(id)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to view statistics for this assignment":
			c.Fail(403, "无权限查看此作业统计信息")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	c.Success(stats)
}