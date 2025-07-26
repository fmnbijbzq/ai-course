package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
)

// SubmissionRepository 提交仓储接口
type SubmissionRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, submission *model.Submission) error
	GetByID(ctx context.Context, id uint) (*model.Submission, error)
	Update(ctx context.Context, submission *model.Submission) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentID uint) (*model.Submission, error)
	GetByAssignmentID(ctx context.Context, assignmentID uint, offset, limit int) ([]*model.Submission, int64, error)
	GetByStudentID(ctx context.Context, studentID uint, offset, limit int) ([]*model.Submission, int64, error)
	
	// 提交详情（包含答案）
	GetDetailByID(ctx context.Context, id uint) (*model.Submission, error)
	GetByIDWithDetail(ctx context.Context, id uint) (*model.SubmissionDetail, error)
	GetByAssignmentIDWithPagination(ctx context.Context, assignmentID uint, page, pageSize int, status string) ([]*model.SubmissionDetail, int64, error)
	
	// 统计操作
	CountByAssignmentAndStatus(ctx context.Context, assignmentID uint, status model.SubmissionStatus) (int64, error)
	GetSubmissionStats(ctx context.Context, assignmentID uint) (map[model.SubmissionStatus]int64, error)
	GetStatistics(ctx context.Context, assignmentID uint) (*model.SubmissionStatistics, error)
}

// submissionRepository 提交仓储实现
type submissionRepository struct {
	db    DB
	cache Cache
}

// NewSubmissionRepository 创建提交仓储实例
func NewSubmissionRepository(db DB, cache Cache) SubmissionRepository {
	return &submissionRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建提交
func (r *submissionRepository) Create(ctx context.Context, submission *model.Submission) error {
	if err := r.db.WithContext(ctx).Create(submission).Error; err != nil {
		return fmt.Errorf("create submission failed: %w", err)
	}
	return nil
}

// GetByID 根据ID获取提交
func (r *submissionRepository) GetByID(ctx context.Context, id uint) (*model.Submission, error) {
	var submission model.Submission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&submission).Error
	if err != nil {
		return nil, fmt.Errorf("get submission by id failed: %w", err)
	}
	return &submission, nil
}

// Update 更新提交
func (r *submissionRepository) Update(ctx context.Context, submission *model.Submission) error {
	if err := r.db.WithContext(ctx).Save(submission).Error; err != nil {
		return fmt.Errorf("update submission failed: %w", err)
	}
	return nil
}

// Delete 删除提交
func (r *submissionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Submission{}, id).Error; err != nil {
		return fmt.Errorf("delete submission failed: %w", err)
	}
	return nil
}

// GetByAssignmentAndStudent 根据作业ID和学生ID获取提交
func (r *submissionRepository) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentID uint) (*model.Submission, error) {
	var submission model.Submission
	err := r.db.WithContext(ctx).
		Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).
		First(&submission).Error
	
	if err != nil {
		return nil, fmt.Errorf("get submission by assignment and student failed: %w", err)
	}
	
	return &submission, nil
}

// GetByAssignmentID 根据作业ID获取提交列表
func (r *submissionRepository) GetByAssignmentID(ctx context.Context, assignmentID uint, offset, limit int) ([]*model.Submission, int64, error) {
	db := r.db.WithContext(ctx).Where("assignment_id = ?", assignmentID)
	
	// 获取总数
	var total int64
	if err := db.Model(&model.Submission{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count submissions failed: %w", err)
	}
	
	// 获取列表
	var submissions []*model.Submission
	err := db.Preload("Student").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&submissions).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("get submissions by assignment id failed: %w", err)
	}
	
	return submissions, total, nil
}

// GetByStudentID 根据学生ID获取提交列表
func (r *submissionRepository) GetByStudentID(ctx context.Context, studentID uint, offset, limit int) ([]*model.Submission, int64, error) {
	db := r.db.WithContext(ctx).Where("student_id = ?", studentID)
	
	// 获取总数
	var total int64
	if err := db.Model(&model.Submission{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count submissions failed: %w", err)
	}
	
	// 获取列表
	var submissions []*model.Submission
	err := db.Preload("Assignment").
		Preload("Assignment.Class").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&submissions).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("get submissions by student id failed: %w", err)
	}
	
	return submissions, total, nil
}

// GetDetailByID 获取提交详情（包含答案）
func (r *submissionRepository) GetDetailByID(ctx context.Context, id uint) (*model.Submission, error) {
	var submission model.Submission
	err := r.db.WithContext(ctx).
		Preload("Assignment").
		Preload("Assignment.Questions", func(db DB) DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Student").
		Preload("Answers").
		Preload("Answers.Question").
		Where("id = ?", id).
		First(&submission).Error
	
	if err != nil {
		return nil, fmt.Errorf("get submission detail failed: %w", err)
	}
	
	return &submission, nil
}

// CountByAssignmentAndStatus 根据作业ID和状态统计提交数量
func (r *submissionRepository) CountByAssignmentAndStatus(ctx context.Context, assignmentID uint, status model.SubmissionStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, status).
		Count(&count).Error
	
	if err != nil {
		return 0, fmt.Errorf("count submissions by status failed: %w", err)
	}
	
	return count, nil
}

// GetSubmissionStats 获取作业提交统计信息
func (r *submissionRepository) GetSubmissionStats(ctx context.Context, assignmentID uint) (map[model.SubmissionStatus]int64, error) {
	stats := make(map[model.SubmissionStatus]int64)
	
	// 统计各种状态的提交数量
	statuses := []model.SubmissionStatus{
		model.SubmissionStatusDraft,
		model.SubmissionStatusSubmitted,
		model.SubmissionStatusGraded,
	}
	
	for _, status := range statuses {
		count, err := r.CountByAssignmentAndStatus(ctx, assignmentID, status)
		if err != nil {
			return nil, fmt.Errorf("get submission stats failed: %w", err)
		}
		stats[status] = count
	}
	
	return stats, nil
}

// GetByIDWithDetail 获取提交详情（包含学生信息和答案）
func (r *submissionRepository) GetByIDWithDetail(ctx context.Context, id uint) (*model.SubmissionDetail, error) {
	var submission model.Submission
	err := r.db.WithContext(ctx).
		Preload("Assignment").
		Preload("Assignment.Questions", func(db DB) DB {
			return db.Order("\"order\" ASC")
		}).
		Preload("Student").
		Preload("Answers").
		Preload("Answers.Question").
		Where("id = ?", id).
		First(&submission).Error
	
	if err != nil {
		return nil, fmt.Errorf("get submission detail failed: %w", err)
	}
	
	detail := &model.SubmissionDetail{
		Submission: submission,
		Student:    submission.Student,
		Answers:    submission.Answers,
	}
	
	return detail, nil
}

// GetByAssignmentIDWithPagination 分页获取作业的提交列表（用于批改）
func (r *submissionRepository) GetByAssignmentIDWithPagination(ctx context.Context, assignmentID uint, page, pageSize int, status string) ([]*model.SubmissionDetail, int64, error) {
	db := r.db.WithContext(ctx).Where("assignment_id = ?", assignmentID)
	
	// 如果指定状态，添加状态过滤
	if status != "" {
		db = db.Where("status = ?", status)
	}
	
	// 获取总数
	var total int64
	if err := db.Model(&model.Submission{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count submissions failed: %w", err)
	}
	
	// 计算偏移量
	offset := (page - 1) * pageSize
	
	// 获取提交列表
	var submissions []*model.Submission
	err := db.Preload("Student").
		Preload("Answers").
		Preload("Answers.Question").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&submissions).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("get submissions by assignment id failed: %w", err)
	}
	
	// 转换为 SubmissionDetail
	details := make([]*model.SubmissionDetail, len(submissions))
	for i, submission := range submissions {
		details[i] = &model.SubmissionDetail{
			Submission: *submission,
			Student:    submission.Student,
			Answers:    submission.Answers,
		}
	}
	
	return details, total, nil
}

// GetStatistics 获取作业提交统计
func (r *submissionRepository) GetStatistics(ctx context.Context, assignmentID uint) (*model.SubmissionStatistics, error) {
	stats := &model.SubmissionStatistics{}
	
	// 统计总提交数
	err := r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ?", assignmentID).
		Count(&stats.TotalSubmissions).Error
	if err != nil {
		return nil, fmt.Errorf("count total submissions failed: %w", err)
	}
	
	// 统计草稿数
	err = r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, model.SubmissionStatusDraft).
		Count(&stats.DraftSubmissions).Error
	if err != nil {
		return nil, fmt.Errorf("count draft submissions failed: %w", err)
	}
	
	// 统计已提交数
	err = r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, model.SubmissionStatusSubmitted).
		Count(&stats.SubmittedSubmissions).Error
	if err != nil {
		return nil, fmt.Errorf("count submitted submissions failed: %w", err)
	}
	
	// 统计已批改数
	err = r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, model.SubmissionStatusGraded).
		Count(&stats.GradedSubmissions).Error
	if err != nil {
		return nil, fmt.Errorf("count graded submissions failed: %w", err)
	}
	
	return stats, nil
}