package controller

import (
	"ai-course/internal/base/controller"
	"ai-course/internal/logger"
	"ai-course/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AttachmentController 附件控制器
type AttachmentController struct {
	controller.BaseController
	attachmentService service.AttachmentService
}

// NewAttachmentController 创建附件控制器
func NewAttachmentController(attachmentService service.AttachmentService) *AttachmentController {
	return &AttachmentController{
		attachmentService: attachmentService,
	}
}

// Upload godoc
// @Summary 上传附件
// @Description 教师上传作业附件
// @Tags 附件管理
// @Accept multipart/form-data
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Param file formData file true "文件"
// @Success 200 {object} response.Response "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "作业不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/attachment/assignment/{assignment_id} [post]
func (c *AttachmentController) Upload(ctx *gin.Context) {
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
		logger.Logger.Warn("User not authenticated for file upload")
		c.Unauthorized("用户未认证")
		return
	}

	uploaderID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for file upload")
		c.Unauthorized("用户ID格式无效")
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Logger.Warn("No file uploaded",
			zap.Error(err),
		)
		c.ParamError("请选择要上传的文件")
		return
	}

	// 上传文件
	attachment, err := c.attachmentService.UploadFile(ctx.Request.Context(), file, uint(assignmentID), uploaderID)
	if err != nil {
		logger.Logger.Error("Failed to upload file",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
			zap.Uint("uploader_id", uploaderID),
			zap.String("filename", file.Filename),
		)

		switch err.Error() {
		case "assignment not found":
			c.Fail(404, "作业不存在")
		case "no permission to upload file to this assignment":
			c.Fail(403, "无权限上传文件到此作业")
		case "file type not allowed":
			c.ParamError("不支持的文件类型")
		case "file size exceeds limit (10MB)":
			c.ParamError("文件大小超过限制(10MB)")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("File uploaded successfully",
		zap.Uint("attachment_id", attachment.ID),
		zap.String("filename", file.Filename),
		zap.Uint("assignment_id", uint(assignmentID)),
	)

	c.SuccessWithMessage("文件上传成功", attachment)
}

// GetByAssignment godoc
// @Summary 获取作业附件列表
// @Description 获取指定作业的附件列表
// @Tags 附件管理
// @Produce json
// @Param assignment_id path int true "作业ID"
// @Success 200 {object} response.Response "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/attachment/assignment/{assignment_id} [get]
func (c *AttachmentController) GetByAssignment(ctx *gin.Context) {
	c.InitHandler(ctx)
	assignmentIDStr := ctx.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.ParamError("作业ID格式无效")
		return
	}

	attachments, err := c.attachmentService.GetByAssignmentID(ctx.Request.Context(), uint(assignmentID))
	if err != nil {
		logger.Logger.Error("Failed to get attachments by assignment ID",
			zap.Error(err),
			zap.Uint("assignment_id", uint(assignmentID)),
		)
		c.ServerError(err.Error())
		return
	}

	c.Success(attachments)
}

// Download godoc
// @Summary 下载附件
// @Description 下载指定附件
// @Tags 附件管理
// @Produce application/octet-stream
// @Param id path int true "附件ID"
// @Success 200 {file} file "文件内容"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "附件不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/attachment/{id}/download [get]
func (c *AttachmentController) Download(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("附件ID格式无效")
		return
	}

	// 获取附件详情
	attachment, err := c.attachmentService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		logger.Logger.Error("Failed to get attachment for download",
			zap.Error(err),
			zap.Uint("attachment_id", uint(id)),
		)
		c.Fail(404, "附件不存在")
		return
	}

	// 获取文件路径
	filePath, err := c.attachmentService.DownloadFile(ctx.Request.Context(), uint(id))
	if err != nil {
		logger.Logger.Error("Failed to get file path for download",
			zap.Error(err),
			zap.Uint("attachment_id", uint(id)),
		)

		switch err.Error() {
		case "attachment not found":
			c.Fail(404, "附件不存在")
		case "file not found on disk":
			c.Fail(404, "文件不存在")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+attachment.FileName)
	ctx.Header("Content-Type", attachment.ContentType)

	// 发送文件
	ctx.File(filePath)

	logger.Logger.Info("File downloaded",
		zap.Uint("attachment_id", uint(id)),
		zap.String("filename", attachment.FileName),
	)
}

// Delete godoc
// @Summary 删除附件
// @Description 删除指定附件
// @Tags 附件管理
// @Produce json
// @Param id path int true "附件ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "附件不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/attachment/{id} [delete]
func (c *AttachmentController) Delete(ctx *gin.Context) {
	c.InitHandler(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.ParamError("附件ID格式无效")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Logger.Warn("User not authenticated for file deletion")
		c.Unauthorized("用户未认证")
		return
	}

	uploaderID, ok := userID.(uint)
	if !ok {
		logger.Logger.Warn("Invalid user ID format for file deletion")
		c.Unauthorized("用户ID格式无效")
		return
	}

	// 删除文件
	err = c.attachmentService.DeleteFile(ctx.Request.Context(), uint(id), uploaderID)
	if err != nil {
		logger.Logger.Error("Failed to delete attachment",
			zap.Error(err),
			zap.Uint("attachment_id", uint(id)),
			zap.Uint("user_id", uploaderID),
		)

		switch err.Error() {
		case "attachment not found":
			c.Fail(404, "附件不存在")
		case "no permission to delete this attachment":
			c.Fail(403, "无权限删除此附件")
		default:
			c.ServerError(err.Error())
		}
		return
	}

	logger.Logger.Info("Attachment deleted successfully",
		zap.Uint("attachment_id", uint(id)),
		zap.Uint("user_id", uploaderID),
	)

	c.SuccessWithMessage("附件删除成功", nil)
}