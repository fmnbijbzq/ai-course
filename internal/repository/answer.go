package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
)

// AnswerRepository 答案仓储接口
type AnswerRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, answer *model.Answer) error
	GetByID(ctx context.Context, id uint) (*model.Answer, error)
	Update(ctx context.Context, answer *model.Answer) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	GetBySubmissionID(ctx context.Context, submissionID uint) ([]*model.Answer, error)
	GetByQuestionID(ctx context.Context, questionID uint) ([]*model.Answer, error)
	GetBySubmissionAndQuestion(ctx context.Context, submissionID, questionID uint) (*model.Answer, error)
	
	// 批量操作
	CreateBatch(ctx context.Context, answers []*model.Answer) error
	UpdateBatch(ctx context.Context, answers []*model.Answer) error
	DeleteBySubmissionID(ctx context.Context, submissionID uint) error
}

// answerRepository 答案仓储实现
type answerRepository struct {
	db    DB
	cache Cache
}

// NewAnswerRepository 创建答案仓储实例
func NewAnswerRepository(db DB, cache Cache) AnswerRepository {
	return &answerRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建答案
func (r *answerRepository) Create(ctx context.Context, answer *model.Answer) error {
	if err := r.db.WithContext(ctx).Create(answer).Error; err != nil {
		return fmt.Errorf("create answer failed: %w", err)
	}
	return nil
}

// GetByID 根据ID获取答案
func (r *answerRepository) GetByID(ctx context.Context, id uint) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&answer).Error
	if err != nil {
		return nil, fmt.Errorf("get answer by id failed: %w", err)
	}
	return &answer, nil
}

// Update 更新答案
func (r *answerRepository) Update(ctx context.Context, answer *model.Answer) error {
	if err := r.db.WithContext(ctx).Save(answer).Error; err != nil {
		return fmt.Errorf("update answer failed: %w", err)
	}
	return nil
}

// Delete 删除答案
func (r *answerRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Answer{}, id).Error; err != nil {
		return fmt.Errorf("delete answer failed: %w", err)
	}
	return nil
}

// GetBySubmissionID 根据提交ID获取答案列表
func (r *answerRepository) GetBySubmissionID(ctx context.Context, submissionID uint) ([]*model.Answer, error) {
	var answers []*model.Answer
	err := r.db.WithContext(ctx).
		Preload("Question").
		Where("submission_id = ?", submissionID).
		Order("question_id ASC").
		Find(&answers).Error
	
	if err != nil {
		return nil, fmt.Errorf("get answers by submission id failed: %w", err)
	}
	
	return answers, nil
}

// GetByQuestionID 根据题目ID获取答案列表
func (r *answerRepository) GetByQuestionID(ctx context.Context, questionID uint) ([]*model.Answer, error) {
	var answers []*model.Answer
	err := r.db.WithContext(ctx).
		Preload("Submission").
		Preload("Submission.Student").
		Where("question_id = ?", questionID).
		Find(&answers).Error
	
	if err != nil {
		return nil, fmt.Errorf("get answers by question id failed: %w", err)
	}
	
	return answers, nil
}

// GetBySubmissionAndQuestion 根据提交ID和题目ID获取答案
func (r *answerRepository) GetBySubmissionAndQuestion(ctx context.Context, submissionID, questionID uint) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.WithContext(ctx).
		Where("submission_id = ? AND question_id = ?", submissionID, questionID).
		First(&answer).Error
	
	if err != nil {
		return nil, fmt.Errorf("get answer by submission and question failed: %w", err)
	}
	
	return &answer, nil
}

// CreateBatch 批量创建答案
func (r *answerRepository) CreateBatch(ctx context.Context, answers []*model.Answer) error {
	if len(answers) == 0 {
		return nil
	}
	
	// 使用事务批量创建
	return r.db.WithContext(ctx).Transaction(func(tx DB) error {
		for _, answer := range answers {
			if err := tx.Create(answer).Error; err != nil {
				return fmt.Errorf("create answer batch failed: %w", err)
			}
		}
		return nil
	})
}

// UpdateBatch 批量更新答案
func (r *answerRepository) UpdateBatch(ctx context.Context, answers []*model.Answer) error {
	if len(answers) == 0 {
		return nil
	}
	
	// 使用事务批量更新
	return r.db.WithContext(ctx).Transaction(func(tx DB) error {
		for _, answer := range answers {
			if err := tx.Save(answer).Error; err != nil {
				return fmt.Errorf("update answer batch failed: %w", err)
			}
		}
		return nil
	})
}

// DeleteBySubmissionID 根据提交ID删除所有答案
func (r *answerRepository) DeleteBySubmissionID(ctx context.Context, submissionID uint) error {
	err := r.db.WithContext(ctx).
		Where("submission_id = ?", submissionID).
		Delete(&model.Answer{}).Error
	
	if err != nil {
		return fmt.Errorf("delete answers by submission id failed: %w", err)
	}
	
	return nil
}