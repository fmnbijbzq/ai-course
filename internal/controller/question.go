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

// QuestionController 题目控制器
type QuestionController struct {
	controller.BaseController
	questionService   service.QuestionService
	assignmentService service.AssignmentService
}

// NewQuestionController 创建题目控制器
func NewQuestionController(questionService service.QuestionService, assignmentService service.AssignmentService) *QuestionController {
	return &QuestionController{
		questionService:   questionService,
		assignmentService: assignmentService,
	}
}

// Create godoc
// @Summary 为作业添加题目
// @Description 教师为指定作业添加题目
// @Tags 题目管理
// @Accept json
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Param request body model.CreateQuestionRequest true "题目信息"
// @Success 200 {object} response.Response "添加成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{assignment_id}/question [post]
func (c *QuestionController) Create(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	var req model.CreateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid create question request",
			zap.Error(err),
		)
		c.ParamError("创建题目参数无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for create question")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for create question")
		c.Unauthorized("用户ID格式无效")
		return
	}

	// 验证作业权限已在服务层处理，这里直接调用服务
	question, err := c.questionService.CreateQuestion(ctx.Request.Context(), &req, uint(assignmentID), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to create question",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to add question to this assignment":
			c.Fail(403, "无权限为此作业添加题目")
		case "cannot add question to published assignment":
			c.Fail(400, "已发布的作业不能添加题目")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Question created successfully",
		zap.Uint("assignment_id", uint(assignmentID)),
		zap.Uint("question_id", question.ID),
		zap.Uint("teacher_id", teacherID),
		zap.String("question_type", string(req.Type)),
	)

	c.SuccessWithMessage("创建题目成功", question)
}

// Update godoc
// @Summary 更新题目
// @Description 教师更新题目信息
// @Tags 题目管理
// @Accept json
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Param question_id path int true "题目ID"
// @Param request body model.UpdateQuestionRequest true "题目信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "题目不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{assignment_id}/question/{question_id} [put]
func (c *QuestionController) Update(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	questionIDStr := ctx.Param("question_id")
	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.ParamError("题目ID格式无效")
		return
	}

	var req model.UpdateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Logger.Warn("Invalid update question request",
			zap.Error(err),
		)
		c.ParamError("更新题目参数无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for update question")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for update question")
		c.Unauthorized("用户ID格式无效")
		return
	}

	// 更新题目
	question, err := c.questionService.UpdateQuestion(ctx.Request.Context(), uint(questionID), &req, teacherID)
	if err != nil {
		logger.Logger.Error("Failed to update question",
			zap.Error(err),
			zap.Uint("question_id", uint(questionID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "question not found":
			c.Fail(404, "题目不存在")
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to update this question":
			c.Fail(403, "无权限修改此题目")
		case "cannot update question in published assignment":
			c.Fail(400, "已发布的作业不能修改题目")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Question updated successfully",
		zap.Uint("assignment_id", uint(assignmentID)),
		zap.Uint("question_id", question.ID),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("更新题目成功", question)
}

// Delete godoc
// @Summary 删除题目
// @Description 教师删除题目
// @Tags 题目管理
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Param question_id path int true "题目ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "题目不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{assignment_id}/question/{question_id} [delete]
func (c *QuestionController) Delete(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	questionIDStr := ctx.Param("question_id")
	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.ParamError("题目ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for delete question")
		c.Unauthorized("用户未认证")
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for delete question")
		c.Unauthorized("用户ID格式无效")
		return
	}

	// 删除题目
	err = c.questionService.DeleteQuestion(ctx.Request.Context(), uint(questionID), teacherID)
	if err != nil {
		logger.Logger.Error("Failed to delete question",
			zap.Error(err),
			zap.Uint("question_id", uint(questionID)),
			zap.Uint("teacher_id", teacherID),
		)
		
		switch err.Error() {
		case "question not found":
			c.Fail(404, "题目不存在")
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "teacher has no permission to delete this question":
			c.Fail(403, "无权限删除此题目")
		case "cannot delete question from published assignment":
			c.Fail(400, "已发布的作业不能删除题目")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Question deleted successfully",
		zap.Uint("assignment_id", uint(assignmentID)),
		zap.Uint("question_id", uint(questionID)),
		zap.Uint("teacher_id", teacherID),
	)

	c.SuccessWithMessage("删除题目成功", nil)
}

// List godoc
// @Summary 获取作业题目列表
// @Description 获取指定作业的所有题目
// @Tags 题目管理
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/assignment/{assignment_id}/question [get]
func (c *QuestionController) List(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	// 获取作业题目列表
	questions, err := c.questionService.GetQuestionsByAssignmentID(ctx.Request.Context(), uint(assignmentID))
	if err != nil {
		logger.Logger.Error("Failed to get questions for assignment",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
		)
		c.ServerError(err.Error())
		return
	}

	c.Success(gin.H{
		"assignment_id": assignmentID,
		"questions":     questions,
		"total":         len(questions),
	})
}