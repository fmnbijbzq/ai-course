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

// SubmissionController 提交控制器
type SubmissionController struct {
	controller.BaseController
	submissionService service.SubmissionService
}

// NewSubmissionController 创建提交控制器
func NewSubmissionController(submissionService service.SubmissionService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
	}
}

// SaveDraft godoc
// @Summary 保存作业草稿
// @Description 学生保存作业答案为草稿
// @Tags 作业提交
// @Accept json
// @Produce json
// @Param request body model.SubmissionRequest true "提交信息"
// @Success 200 {object} response.Response "保存成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/submission/draft [post]
func (c *SubmissionController) SaveDraft(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req model.SubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid save draft request",
			zap.Error(err),
		)
		c.ParamError("保存草稿参数无效")
		return
	}

	// 强制设置为草稿状态
	req.Status = model.SubmissionStatusDraft

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for save draft")
		c.Unauthorized("用户未认证")
		return
	}

	studentID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for save draft")
		c.Unauthorized("用户ID格式无效")
		return
	}

	submission, err := c.submissionService.CreateOrUpdateSubmission(ctx.Request.Context(), &req, studentID)
	if err != nil {
		logger.Logger.Error("Failed to save draft",
			zap.Error(err),
			zap.Uint("student_id", studentID),
			zap.Uint("assignment_id", req.AssignmentID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "assignment is not published":
			c.Fail(400, "作业未发布")
		case "cannot change submitted assignment back to draft":
			c.Fail(400, "已提交的作业不能修改为草稿")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Draft saved successfully",
		zap.Uint("submission_id", submission.ID),
		zap.Uint("student_id", studentID),
		zap.Uint("assignment_id", req.AssignmentID),
	)

	c.SuccessWithMessage("保存草稿成功", submission)
}

// Submit godoc
// @Summary 提交作业
// @Description 学生提交作业
// @Tags 作业提交
// @Accept json
// @Produce json
// @Param request body model.SubmissionRequest true "提交信息"
// @Success 200 {object} response.Response "提交成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/submission/submit [post]
func (c *SubmissionController) Submit(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req model.SubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid submit request",
			zap.Error(err),
		)
		c.ParamError("提交作业参数无效")
		return
	}

	// 强制设置为提交状态
	req.Status = model.SubmissionStatusSubmitted

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for submit")
		c.Unauthorized("用户未认证")
		return
	}

	studentID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for submit")
		c.Unauthorized("用户ID格式无效")
		return
	}

	submission, err := c.submissionService.CreateOrUpdateSubmission(ctx.Request.Context(), &req, studentID)
	if err != nil {
		logger.Logger.Error("Failed to submit assignment",
			zap.Error(err),
			zap.Uint("student_id", studentID),
			zap.Uint("assignment_id", req.AssignmentID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "assignment is not published":
			c.Fail(400, "作业未发布")
		case "assignment deadline has passed":
			c.Fail(400, "作业已过截止时间")
		case "assignment already submitted":
			c.Fail(400, "作业已提交")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Assignment submitted successfully",
		zap.Uint("submission_id", submission.ID),
		zap.Uint("student_id", studentID),
		zap.Uint("assignment_id", req.AssignmentID),
	)

	c.SuccessWithMessage("提交作业成功", submission)
}

// GetStudentAssignments godoc
// @Summary 获取学生作业列表
// @Description 获取当前学生的所有作业列表
// @Tags 作业提交
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/submission/student/assignments [get]
func (c *SubmissionController) GetStudentAssignments(ctx *gin.Context) {
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
		logger.Logger.Warn("User not authenticated for get student assignments")
		c.Unauthorized("用户未认证")
		return
	}

	studentID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for get student assignments")
		c.Unauthorized("用户ID格式无效")
		return
	}

	assignments, total, err := c.submissionService.GetStudentSubmissions(ctx.Request.Context(), studentID, page, pageSize)
	if err != nil {
		logger.Logger.Error("Failed to get student assignments",
			zap.Error(err),
			zap.Uint("student_id", studentID),
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

// GetAssignmentForStudent godoc
// @Summary 获取学生特定作业详情
// @Description 获取学生特定作业的详细信息和提交状态
// @Tags 作业提交
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/submission/assignment/{assignment_id} [get]
func (c *SubmissionController) GetAssignmentForStudent(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for get assignment")
		c.Unauthorized("用户未认证")
		return
	}

	studentID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for get assignment")
		c.Unauthorized("用户ID格式无效")
		return
	}

	assignment, err := c.submissionService.GetStudentSubmissionByAssignment(ctx.Request.Context(), uint(assignmentID), studentID)
	if err != nil {
		logger.Logger.Error("Failed to get assignment for student",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
			zap.Uint("student_id", studentID),
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

// GetSubmissionDetail godoc
// @Summary 获取提交详情
// @Description 获取提交的详细信息
// @Tags 作业提交
// @Produce json
// @Param id path int true "提交ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "提交不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/submission/{id} [get]
func (c *SubmissionController) GetSubmissionDetail(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("提交ID格式无效")
		return
	}

	// 获取当前用户ID（用于权限验证）
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for get submission detail")
		c.Unauthorized("用户未认证")
		return
	}

	submission, err := c.submissionService.GetSubmissionDetail(ctx.Request.Context(), uint(id))
	if err != nil {
		logger.Logger.Error("Failed to get submission detail",
			zap.Error(err),
			zap.Uint("submission_id", uint(id)),
		)
		
		if err.Error() == "submission not found" {
			c.Fail(404, "提交不存在")
			return
		}
		c.ServerError(err.Error())
		return
	}

	// 简单的权限检查：只有提交者本人或任课教师可以查看
	currentUserID := userID.(uint)
	if submission.StudentID != currentUserID && submission.Assignment.TeacherID != currentUserID {
		c.Fail(403, "无权限查看此提交")
		return
	}

	c.Success(submission)
}