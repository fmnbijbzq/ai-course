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

// GradingController 批改控制器
type GradingController struct {
	controller.BaseController
	submissionService service.SubmissionService
	gradingService    service.GradingService
}

// NewGradingController 创建批改控制器
func NewGradingController(submissionService service.SubmissionService, gradingService service.GradingService) *GradingController {
	return &GradingController{
		submissionService: submissionService,
		gradingService:    gradingService,
	}
}

// GetSubmissionsForGrading godoc
// @Summary 获取待批改提交列表
// @Description 教师获取指定作业的待批改提交列表
// @Tags 作业批改
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "提交状态过滤" Enums(submitted,graded)
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/grading/assignment/{assignment_id}/submissions [get]
func (c *GradingController) GetSubmissionsForGrading(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	status := ctx.Query("status")
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for get submissions")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for get submissions")
		c.Unauthorized("用户ID格式无效")
		return
	}

	submissions, total, err := c.gradingService.GetSubmissionsForGrading(ctx.Request.Context(), uint(assignmentID), teacherID, page, pageSize, status)
	if err != nil {
		logger.Logger.Error("Failed to get submissions for grading",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to grade this assignment":
			c.Fail(403, "无权限批改此作业")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	response := gin.H{
		"list":  submissions,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}

	c.Success(response)
}

// GetGradingDetail godoc
// @Summary 获取批改详情
// @Description 教师获取特定提交的批改详情
// @Tags 作业批改
// @Produce json
// @Param submission_id path int true "提交ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "提交不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/grading/submission/{submission_id} [get]
func (c *GradingController) GetGradingDetail(ctx *gin.Context) {
	c.InitHandler(ctx)
	submissionIDStr := ctx.Param("submission_id")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
	if err != nil {
		c.ParamError("提交ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for get grading detail")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for get grading detail")
		c.Unauthorized("用户ID格式无效")
		return
	}

	detail, err := c.submissionService.GetGradingDetail(ctx.Request.Context(), uint(submissionID), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to get grading detail",
			zap.Error(err),
			zap.Uint("submission_id", uint(submissionID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "submission not found":
			c.Fail(404, "提交不存在")
		case "teacher has no permission to grade this submission":
			c.Fail(403, "无权限批改此提交")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	c.Success(detail)
}

// GradeSubmission godoc
// @Summary 批改提交
// @Description 教师批改学生提交的作业
// @Tags 作业批改
// @Accept json
// @Produce json
// @Param submission_id path int true "提交ID"
// @Param request body model.GradeSubmissionRequest true "批改信息"
// @Success 200 {object} response.Response "批改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "提交不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/grading/submission/{submission_id} [post]
func (c *GradingController) GradeSubmission(ctx *gin.Context) {
	c.InitHandler(ctx)
	submissionIDStr := ctx.Param("submission_id")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
	if err != nil {
		c.ParamError("提交ID格式无效")
		return
	}

	var req model.GradeSubmissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid grade submission request",
			zap.Error(err),
		)
		c.ParamError("批改参数无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for grade submission")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for grade submission")
		c.Unauthorized("用户ID格式无效")
		return
	}

	submission, err := c.gradingService.GradeSubmission(ctx.Request.Context(), uint(submissionID), &req, teacherID)
	if err != nil {
		logger.Logger.Error("Failed to grade submission",
			zap.Error(err),
			zap.Uint("submission_id", uint(submissionID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "submission not found":
			c.Fail(404, "提交不存在")
		case "teacher has no permission to grade this submission":
			c.Fail(403, "无权限批改此提交")
		case "submission is not submitted yet":
			c.Fail(400, "作业尚未提交")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Submission graded successfully",
		zap.Uint("submission_id", uint(submissionID)),
		zap.Uint("teacher_id", teacherID),
		zap.Int("score", submission.Score),
	)

	c.SuccessWithMessage("批改完成", submission)
}

// BatchGrade godoc
// @Summary 批量快速批改
// @Description 教师对多个提交进行批量快速批改
// @Tags 作业批改
// @Accept json
// @Produce json
// @Param request body model.BatchGradeRequest true "批量批改信息"
// @Success 200 {object} response.Response "批改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/grading/batch [post]
func (c *GradingController) BatchGrade(ctx *gin.Context) {
	c.InitHandler(ctx)
	var req model.BatchGradeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid batch grade request",
			zap.Error(err),
		)
		c.ParamError("批量批改参数无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for batch grade")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for batch grade")
		c.Unauthorized("用户ID格式无效")
		return
	}

	results, err := c.gradingService.BatchGrade(ctx.Request.Context(), &req, teacherID)
	if err != nil {
		logger.Logger.Error("Failed to batch grade",
			zap.Error(err),
			zap.Uint("teacher_id", teacherID),
		)
		c.ServerError(err.Error())
		return
	}

	logger.Logger.Info("Batch grading completed",
		zap.Uint("teacher_id", teacherID),
		zap.Int("total_processed", len(results)),
	)

	c.SuccessWithMessage("批量批改完成", gin.H{
		"results": results,
		"total":   len(results),
	})
}

// PublishGrades godoc
// @Summary 发布成绩
// @Description 教师发布作业成绩，学生可查看批改结果
// @Tags 作业批改
// @Accept json
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Success 200 {object} response.Response "发布成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/grading/assignment/{assignment_id}/publish [post]
func (c *GradingController) PublishGrades(ctx *gin.Context) {
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
		logger.Logger.Warn("User not authenticated for publish grades")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for publish grades")
		c.Unauthorized("用户ID格式无效")
		return
	}

	err = c.gradingService.PublishGrades(ctx.Request.Context(), uint(assignmentID), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to publish grades",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to publish grades for this assignment":
			c.Fail(403, "无权限发布此作业成绩")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Grades published successfully",
		zap.Uint("assignment_id", uint(assignmentID)),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("成绩发布成功", nil)
}

// GetGradingProgress godoc
// @Summary 获取批改进度
// @Description 获取作业的批改进度统计
// @Tags 作业批改
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/grading/assignment/{assignment_id}/progress [get]
func (c *GradingController) GetGradingProgress(ctx *gin.Context) {
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
		logger.Logger.Warn("User not authenticated for get grading progress")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for get grading progress")
		c.Unauthorized("用户ID格式无效")
		return
	}

	progress, err := c.gradingService.GetGradingProgress(ctx.Request.Context(), uint(assignmentID), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to get grading progress",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to view grading progress for this assignment":
			c.Fail(403, "无权限查看此作业批改进度")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	c.Success(progress)
}