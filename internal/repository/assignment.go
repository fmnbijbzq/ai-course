package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
)

// AssignmentRepository 作业仓储接口
type AssignmentRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, assignment *model.Assignment) error
	GetByID(ctx context.Context, id uint) (*model.Assignment, error)
	Update(ctx context.Context, assignment *model.Assignment) error
	Delete(ctx context.Context, id uint) error
	
	// 查询操作
	GetByTeacherID(ctx context.Context, teacherID uint, offset, limit int) ([]*model.Assignment, int64, error)
	GetByClassID(ctx context.Context, classID uint, offset, limit int) ([]*model.Assignment, int64, error)
	GetByStudentID(ctx context.Context, studentID uint, offset, limit int) ([]*model.Assignment, int64, error)
	GetPublishedByClassID(ctx context.Context, classID uint) ([]*model.Assignment, error)
	
	// 作业详情（包含题目和附件）
	GetDetailByID(ctx context.Context, id uint) (*model.Assignment, error)
	
	// 统计操作
	GetSubmissionStats(ctx context.Context, assignmentID uint) (*model.AssignmentStatistics, error)
}

// assignmentRepository 作业仓储实现
type assignmentRepository struct {
	db    DB
	cache Cache
}

// NewAssignmentRepository 创建作业仓储实例
func NewAssignmentRepository(db DB, cache Cache) AssignmentRepository {
	return &assignmentRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建作业
func (r *assignmentRepository) Create(ctx context.Context, assignment *model.Assignment) error {
	if err := r.db.WithContext(ctx).Create(assignment).Error; err != nil {
		return fmt.Errorf("create assignment failed: %w", err)
	}
	return nil
}

// GetByID 根据ID获取作业
func (r *assignmentRepository) GetByID(ctx context.Context, id uint) (*model.Assignment, error) {
	var assignment model.Assignment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&assignment).Error
	if err != nil {
		return nil, fmt.Errorf("get assignment by id failed: %w", err)
	}
	return &assignment, nil
}

// Update 更新作业
func (r *assignmentRepository) Update(ctx context.Context, assignment *model.Assignment) error {
	if err := r.db.WithContext(ctx).Save(assignment).Error; err != nil {
		return fmt.Errorf("update assignment failed: %w", err)
	}
	return nil
}

// Delete 删除作业
func (r *assignmentRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Assignment{}, id).Error; err != nil {
		return fmt.Errorf("delete assignment failed: %w", err)
	}
	return nil
}

// GetByTeacherID 根据教师ID获取作业列表
func (r *assignmentRepository) GetByTeacherID(ctx context.Context, teacherID uint, offset, limit int) ([]*model.Assignment, int64, error) {
	db := r.db.WithContext(ctx).Where("teacher_id = ?", teacherID)
	
	// 获取总数
	var total int64
	if err := db.Model(&model.Assignment{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count assignments failed: %w", err)
	}
	
	// 获取列表
	var assignments []*model.Assignment
	err := db.Preload("Class").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&assignments).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("get assignments by teacher id failed: %w", err)
	}
	
	return assignments, total, nil
}

// GetByClassID 根据班级ID获取作业列表
func (r *assignmentRepository) GetByClassID(ctx context.Context, classID uint, offset, limit int) ([]*model.Assignment, int64, error) {
	db := r.db.WithContext(ctx).Where("class_id = ?", classID)
	
	// 获取总数
	var total int64
	if err := db.Model(&model.Assignment{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count assignments failed: %w", err)
	}
	
	// 获取列表
	var assignments []*model.Assignment
	err := db.Preload("Teacher").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&assignments).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("get assignments by class id failed: %w", err)
	}
	
	return assignments, total, nil
}

// GetByStudentID 根据学生ID获取作业列表（通过班级关联）
func (r *assignmentRepository) GetByStudentID(ctx context.Context, studentID uint, offset, limit int) ([]*model.Assignment, int64, error) {
	// 首先获取学生所在的班级
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", studentID).First(&user).Error; err != nil {
		return nil, 0, fmt.Errorf("get student failed: %w", err)
	}
	
	// TODO: 这里需要根据实际的学生-班级关联关系来实现
	// 目前假设通过 User 表的某个字段关联班级，需要根据实际模型调整
	
	var assignments []*model.Assignment
	var total int64
	
	// 暂时返回空结果，实际实现需要根据学生-班级关联关系查询
	return assignments, total, nil
}

// GetPublishedByClassID 获取班级的已发布作业
func (r *assignmentRepository) GetPublishedByClassID(ctx context.Context, classID uint) ([]*model.Assignment, error) {
	var assignments []*model.Assignment
	err := r.db.WithContext(ctx).
		Where("class_id = ? AND status = ?", classID, "published").
		Order("created_at DESC").
		Find(&assignments).Error
	
	if err != nil {
		return nil, fmt.Errorf("get published assignments failed: %w", err)
	}
	
	return assignments, nil
}

// GetDetailByID 获取作业详情（包含题目和附件）
func (r *assignmentRepository) GetDetailByID(ctx context.Context, id uint) (*model.Assignment, error) {
	var assignment model.Assignment
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("Teacher").
		Preload("Questions", func(db DB) DB {
			return db.Where("").Order("\"order\" ASC") // 按题目顺序排序
		}).
		Preload("Attachments").
		Where("id = ?", id).
		First(&assignment).Error
	
	if err != nil {
		return nil, fmt.Errorf("get assignment detail failed: %w", err)
	}
	
	return &assignment, nil
}

// GetSubmissionStats 获取作业提交统计信息
func (r *assignmentRepository) GetSubmissionStats(ctx context.Context, assignmentID uint) (*model.AssignmentStatistics, error) {
	stats := &model.AssignmentStatistics{}
	
	// 获取该作业所属班级的学生总数
	// TODO: 需要根据实际的学生-班级关联关系来实现
	// 这里暂时写一个示例实现
	
	// 获取已提交的学生数
	var submittedCount int64
	if err := r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status IN ?", assignmentID, []string{"submitted", "graded"}).
		Count(&submittedCount).Error; err != nil {
		return nil, fmt.Errorf("count submitted assignments failed: %w", err)
	}
	
	// 获取已批改的学生数
	var gradedCount int64
	if err := r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, "graded").
		Count(&gradedCount).Error; err != nil {
		return nil, fmt.Errorf("count graded assignments failed: %w", err)
	}
	
	// 计算平均分
	var avgScore float64
	if err := r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, "graded").
		Select("AVG(score)").
		Scan(&avgScore).Error; err != nil {
		return nil, fmt.Errorf("calculate average score failed: %w", err)
	}
	
	// 获取最高分和最低分
	var maxScore, minScore int
	r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, "graded").
		Select("MAX(score)").
		Scan(&maxScore)
	
	r.db.WithContext(ctx).
		Model(&model.Submission{}).
		Where("assignment_id = ? AND status = ?", assignmentID, "graded").
		Select("MIN(score)").
		Scan(&minScore)
	
	stats.SubmittedCount = int(submittedCount)
	stats.GradedCount = int(gradedCount)
	stats.AverageScore = avgScore
	stats.MaxScore = maxScore
	stats.MinScore = minScore
	
	// TODO: 计算总学生数和提交率
	stats.TotalStudents = 30 // 暂时硬编码，需要根据实际业务逻辑实现
	if stats.TotalStudents > 0 {
		stats.SubmissionRate = float64(stats.SubmittedCount) / float64(stats.TotalStudents) * 100
	}
	
	return stats, nil
}