package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
)

// AttachmentRepository 附件仓储接口
type AttachmentRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, attachment *model.Attachment) error
	GetByID(ctx context.Context, id uint) (*model.Attachment, error)
	Update(ctx context.Context, attachment *model.Attachment) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	GetByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.Attachment, error)
	GetByUploaderID(ctx context.Context, uploaderID uint) ([]*model.Attachment, error)
}

// attachmentRepository 附件仓储实现
type attachmentRepository struct {
	db    DB
	cache Cache
}

// NewAttachmentRepository 创建附件仓储实例
func NewAttachmentRepository(db DB, cache Cache) AttachmentRepository {
	return &attachmentRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建附件
func (r *attachmentRepository) Create(ctx context.Context, attachment *model.Attachment) error {
	if err := r.db.WithContext(ctx).Create(attachment).Error; err != nil {
		return fmt.Errorf("create attachment failed: %w", err)
	}
	return nil
}

// GetByID 根据ID获取附件
func (r *attachmentRepository) GetByID(ctx context.Context, id uint) (*model.Attachment, error) {
	var attachment model.Attachment
	err := r.db.WithContext(ctx).
		Preload("Assignment").
		Preload("Uploader").
		Where("id = ?", id).
		First(&attachment).Error
	
	if err != nil {
		return nil, fmt.Errorf("get attachment by id failed: %w", err)
	}
	
	return &attachment, nil
}

// Update 更新附件
func (r *attachmentRepository) Update(ctx context.Context, attachment *model.Attachment) error {
	if err := r.db.WithContext(ctx).Save(attachment).Error; err != nil {
		return fmt.Errorf("update attachment failed: %w", err)
	}
	return nil
}

// Delete 删除附件
func (r *attachmentRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Attachment{}, id).Error; err != nil {
		return fmt.Errorf("delete attachment failed: %w", err)
	}
	return nil
}

// GetByAssignmentID 根据作业ID获取附件列表
func (r *attachmentRepository) GetByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	err := r.db.WithContext(ctx).
		Preload("Uploader").
		Where("assignment_id = ?", assignmentID).
		Order("created_at DESC").
		Find(&attachments).Error
	
	if err != nil {
		return nil, fmt.Errorf("get attachments by assignment id failed: %w", err)
	}
	
	return attachments, nil
}

// GetByUploaderID 根据上传者ID获取附件列表
func (r *attachmentRepository) GetByUploaderID(ctx context.Context, uploaderID uint) ([]*model.Attachment, error) {
	var attachments []*model.Attachment
	err := r.db.WithContext(ctx).
		Preload("Assignment").
		Where("uploader_id = ?", uploaderID).
		Order("created_at DESC").
		Find(&attachments).Error
	
	if err != nil {
		return nil, fmt.Errorf("get attachments by uploader id failed: %w", err)
	}
	
	return attachments, nil
}