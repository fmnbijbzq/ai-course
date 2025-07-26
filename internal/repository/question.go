package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
)

// QuestionRepository 题目仓储接口
type QuestionRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, question *model.Question) error
	GetByID(ctx context.Context, id uint) (*model.Question, error)
	Update(ctx context.Context, question *model.Question) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	GetByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.Question, error)
	GetByAssignmentIDWithOrder(ctx context.Context, assignmentID uint) ([]*model.Question, error)
	
	// 批量操作
	CreateBatch(ctx context.Context, questions []*model.Question) error
	DeleteByAssignmentID(ctx context.Context, assignmentID uint) error
}

// questionRepository 题目仓储实现
type questionRepository struct {
	db    DB
	cache Cache
}

// NewQuestionRepository 创建题目仓储实例
func NewQuestionRepository(db DB, cache Cache) QuestionRepository {
	return &questionRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建题目
func (r *questionRepository) Create(ctx context.Context, question *model.Question) error {
	if err := r.db.WithContext(ctx).Create(question).Error; err != nil {
		return fmt.Errorf("create question failed: %w", err)
	}
	return nil
}

// GetByID 根据ID获取题目
func (r *questionRepository) GetByID(ctx context.Context, id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, fmt.Errorf("get question by id failed: %w", err)
	}
	return &question, nil
}

// Update 更新题目
func (r *questionRepository) Update(ctx context.Context, question *model.Question) error {
	if err := r.db.WithContext(ctx).Save(question).Error; err != nil {
		return fmt.Errorf("update question failed: %w", err)
	}
	return nil
}

// Delete 删除题目
func (r *questionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Question{}, id).Error; err != nil {
		return fmt.Errorf("delete question failed: %w", err)
	}
	return nil
}

// GetByAssignmentID 根据作业ID获取题目列表
func (r *questionRepository) GetByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.Question, error) {
	var questions []*model.Question
	err := r.db.WithContext(ctx).
		Where("assignment_id = ?", assignmentID).
		Find(&questions).Error
	
	if err != nil {
		return nil, fmt.Errorf("get questions by assignment id failed: %w", err)
	}
	
	return questions, nil
}

// GetByAssignmentIDWithOrder 根据作业ID获取题目列表（按顺序）
func (r *questionRepository) GetByAssignmentIDWithOrder(ctx context.Context, assignmentID uint) ([]*model.Question, error) {
	var questions []*model.Question
	err := r.db.WithContext(ctx).
		Where("assignment_id = ?", assignmentID).
		Order("\"order\" ASC").
		Find(&questions).Error
	
	if err != nil {
		return nil, fmt.Errorf("get questions by assignment id with order failed: %w", err)
	}
	
	return questions, nil
}

// CreateBatch 批量创建题目
func (r *questionRepository) CreateBatch(ctx context.Context, questions []*model.Question) error {
	for _, question := range questions {
		if err := r.Create(ctx, question); err != nil {
			return fmt.Errorf("create question batch failed: %w", err)
		}
	}
	return nil
}

// DeleteByAssignmentID 根据作业ID删除所有题目
func (r *questionRepository) DeleteByAssignmentID(ctx context.Context, assignmentID uint) error {
	err := r.db.WithContext(ctx).
		Where("assignment_id = ?", assignmentID).
		Delete(&model.Question{}).Error
	
	if err != nil {
		return fmt.Errorf("delete questions by assignment id failed: %w", err)
	}
	
	return nil
}