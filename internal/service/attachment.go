package service

import (
	"ai-course/internal/logger"
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// AttachmentService 附件服务接口
type AttachmentService interface {
	// UploadFile 上传文件
	UploadFile(ctx context.Context, file *multipart.FileHeader, assignmentID uint, uploaderID uint) (*model.Attachment, error)
	// GetByAssignmentID 获取作业的附件列表
	GetByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.Attachment, error)
	// GetByID 获取附件详情
	GetByID(ctx context.Context, id uint) (*model.Attachment, error)
	// DownloadFile 下载文件
	DownloadFile(ctx context.Context, id uint) (string, error)
	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, id, userID uint) error
}

// attachmentService 附件服务实现
type attachmentService struct {
	attachmentRepo repository.AttachmentRepository
	assignmentRepo repository.AssignmentRepository
	uploadPath     string
}

// NewAttachmentService 创建附件服务
func NewAttachmentService(
	attachmentRepo repository.AttachmentRepository,
	assignmentRepo repository.AssignmentRepository,
) AttachmentService {
	// 创建上传目录
	uploadPath := "./uploads/attachments"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Logger.Error("Failed to create upload directory", zap.Error(err))
	}

	return &attachmentService{
		attachmentRepo: attachmentRepo,
		assignmentRepo: assignmentRepo,
		uploadPath:     uploadPath,
	}
}

// UploadFile 上传文件
func (s *attachmentService) UploadFile(ctx context.Context, file *multipart.FileHeader, assignmentID uint, uploaderID uint) (*model.Attachment, error) {
	// 验证作业是否存在
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		logger.Logger.Error("Assignment not found for file upload",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return nil, errors.New("assignment not found")
	}

	// 验证用户权限（只有作业创建者可以上传附件）
	if assignment.TeacherID != uploaderID {
		logger.Logger.Warn("User has no permission to upload file to assignment",
			zap.Uint("assignment_id", assignmentID),
			zap.Uint("uploader_id", uploaderID),
			zap.Uint("teacher_id", assignment.TeacherID),
		)
		return nil, errors.New("no permission to upload file to this assignment")
	}

	// 验证文件类型
	if !s.isAllowedFileType(file.Filename) {
		return nil, errors.New("file type not allowed")
	}

	// 验证文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		return nil, errors.New("file size exceeds limit (10MB)")
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%d_%d%s", assignmentID, uploaderID, time.Now().Unix(), ext)
	filePath := filepath.Join(s.uploadPath, fileName)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		logger.Logger.Error("Failed to open uploaded file",
			zap.Error(err),
			zap.String("filename", file.Filename),
		)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		logger.Logger.Error("Failed to create file",
			zap.Error(err),
			zap.String("filepath", filePath),
		)
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		logger.Logger.Error("Failed to save file",
			zap.Error(err),
			zap.String("filepath", filePath),
		)
		// 清理失败的文件
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 创建附件记录
	attachment := &model.Attachment{
		AssignmentID: assignmentID,
		UploaderID:   uploaderID,
		FileName:     file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		ContentType:  file.Header.Get("Content-Type"),
	}

	err = s.attachmentRepo.Create(ctx, attachment)
	if err != nil {
		logger.Logger.Error("Failed to create attachment record",
			zap.Error(err),
			zap.String("filename", file.Filename),
		)
		// 清理文件
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create attachment record: %w", err)
	}

	logger.Logger.Info("File uploaded successfully",
		zap.Uint("attachment_id", attachment.ID),
		zap.String("filename", file.Filename),
		zap.Uint("assignment_id", assignmentID),
	)

	return attachment, nil
}

// GetByAssignmentID 获取作业的附件列表
func (s *attachmentService) GetByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.Attachment, error) {
	attachments, err := s.attachmentRepo.GetByAssignmentID(ctx, assignmentID)
	if err != nil {
		logger.Logger.Error("Failed to get attachments by assignment ID",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return nil, err
	}

	return attachments, nil
}

// GetByID 获取附件详情
func (s *attachmentService) GetByID(ctx context.Context, id uint) (*model.Attachment, error) {
	attachment, err := s.attachmentRepo.GetByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to get attachment by ID",
			zap.Error(err),
			zap.Uint("attachment_id", id),
		)
		return nil, err
	}

	return attachment, nil
}

// DownloadFile 下载文件
func (s *attachmentService) DownloadFile(ctx context.Context, id uint) (string, error) {
	attachment, err := s.attachmentRepo.GetByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Attachment not found for download",
			zap.Error(err),
			zap.Uint("attachment_id", id),
		)
		return "", errors.New("attachment not found")
	}

	// 检查文件是否存在
	if _, err := os.Stat(attachment.FilePath); os.IsNotExist(err) {
		logger.Logger.Error("File not found on disk",
			zap.String("filepath", attachment.FilePath),
			zap.Uint("attachment_id", id),
		)
		return "", errors.New("file not found on disk")
	}

	return attachment.FilePath, nil
}

// DeleteFile 删除文件
func (s *attachmentService) DeleteFile(ctx context.Context, id, userID uint) error {
	attachment, err := s.attachmentRepo.GetByID(ctx, id)
	if err != nil {
		logger.Logger.Error("Attachment not found for deletion",
			zap.Error(err),
			zap.Uint("attachment_id", id),
		)
		return errors.New("attachment not found")
	}

	// 验证权限（只有上传者可以删除）
	if attachment.UploaderID != userID {
		logger.Logger.Warn("User has no permission to delete attachment",
			zap.Uint("attachment_id", id),
			zap.Uint("user_id", userID),
			zap.Uint("uploader_id", attachment.UploaderID),
		)
		return errors.New("no permission to delete this attachment")
	}

	// 删除数据库记录
	err = s.attachmentRepo.Delete(ctx, id)
	if err != nil {
		logger.Logger.Error("Failed to delete attachment record",
			zap.Error(err),
			zap.Uint("attachment_id", id),
		)
		return err
	}

	// 删除文件
	if err := os.Remove(attachment.FilePath); err != nil {
		logger.Logger.Warn("Failed to delete file from disk",
			zap.Error(err),
			zap.String("filepath", attachment.FilePath),
		)
		// 不返回错误，因为数据库记录已删除
	}

	logger.Logger.Info("Attachment deleted successfully",
		zap.Uint("attachment_id", id),
		zap.String("filename", attachment.FileName),
	)

	return nil
}

// isAllowedFileType 检查文件类型是否允许
func (s *attachmentService) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := []string{
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".txt", ".md", ".jpg", ".jpeg", ".png", ".gif", ".zip", ".rar",
	}

	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}

	return false
}