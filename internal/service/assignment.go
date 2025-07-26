package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// AssignmentService 作业服务接口
type AssignmentService interface {
	// 作业管理
	CreateAssignment(ctx context.Context, req *model.CreateAssignmentRequest, teacherID uint) (*model.Assignment, error)
	UpdateAssignment(ctx context.Context, id uint, req *model.UpdateAssignmentRequest, teacherID uint) (*model.Assignment, error)
	DeleteAssignment(ctx context.Context, id uint, teacherID uint) error
	GetAssignment(ctx context.Context, id uint) (*model.Assignment, error)
	GetAssignmentDetail(ctx context.Context, id uint) (*model.AssignmentDetailResponse, error)
	
	// 作业列表
	GetTeacherAssignments(ctx context.Context, teacherID uint, page, pageSize int) ([]*model.AssignmentListResponse, int64, error)
	GetStudentAssignments(ctx context.Context, studentID uint, page, pageSize int) ([]*model.StudentAssignmentResponse, int64, error)
	
	// 发布管理
	PublishAssignment(ctx context.Context, id uint, teacherID uint) error
	UnpublishAssignment(ctx context.Context, id uint, teacherID uint) error
	
	// 统计信息
	GetAssignmentStatistics(ctx context.Context, id uint, teacherID uint) (*model.AssignmentStatistics, error)
}

// assignmentService 作业服务实现
type assignmentService struct {
	assignmentRepo repository.AssignmentRepository
	questionRepo   repository.QuestionRepository
	classRepo      repository.ClassRepository
}

// NewAssignmentService 创建作业服务实例
func NewAssignmentService(
	assignmentRepo repository.AssignmentRepository,
	questionRepo repository.QuestionRepository,
	classRepo repository.ClassRepository,
) AssignmentService {
	return &assignmentService{
		assignmentRepo: assignmentRepo,
		questionRepo:   questionRepo,
		classRepo:      classRepo,
	}
}

// CreateAssignment 创建作业
func (s *assignmentService) CreateAssignment(ctx context.Context, req *model.CreateAssignmentRequest, teacherID uint) (*model.Assignment, error) {
	// 验证班级是否存在且教师有权限
	class, err := s.classRepo.FindByID(ctx, req.ClassID)
	if err != nil {
		return nil, fmt.Errorf("class not found: %w", err)
	}
	
	if class.TeacherID != teacherID {
		return nil, fmt.Errorf("teacher has no permission to create assignment for this class")
	}
	
	// 创建作业
	assignment := &model.Assignment{
		Title:       req.Title,
		Description: req.Description,
		ClassID:     req.ClassID,
		TeacherID:   teacherID,
		Deadline:    req.Deadline,
		TotalScore:  req.TotalScore,
		Status:      req.Status,
	}
	
	if req.Status == "published" {
		now := time.Now()
		assignment.PublishedAt = &now
	}
	
	if err := s.assignmentRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("create assignment failed: %w", err)
	}
	
	// 创建题目
	if len(req.Questions) > 0 {
		for i, questionReq := range req.Questions {
			question := &model.Question{
				AssignmentID:  assignment.ID,
				Type:          questionReq.Type,
				Content:       questionReq.Content,
				Score:         questionReq.Score,
				Order:         questionReq.Order,
				CorrectAnswer: questionReq.CorrectAnswer,
				Reference:     questionReq.Reference,
				Explanation:   questionReq.Explanation,
			}
			
			// 处理选择题选项
			if questionReq.Type == model.QuestionTypeChoice && len(questionReq.Options) > 0 {
				optionsJSON, err := json.Marshal(questionReq.Options)
				if err != nil {
					return nil, fmt.Errorf("marshal question options failed: %w", err)
				}
				question.Options = string(optionsJSON)
			}
			
			if err := s.questionRepo.Create(ctx, question); err != nil {
				return nil, fmt.Errorf("create question %d failed: %w", i+1, err)
			}
		}
	}
	
	return assignment, nil
}

// UpdateAssignment 更新作业
func (s *assignmentService) UpdateAssignment(ctx context.Context, id uint, req *model.UpdateAssignmentRequest, teacherID uint) (*model.Assignment, error) {
	// 获取作业并验证权限
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return nil, fmt.Errorf("teacher has no permission to update this assignment")
	}
	
	// 已发布的作业有限制更新
	if assignment.Status == "published" {
		// 可以延长截止时间，但不能缩短
		if !req.Deadline.IsZero() && req.Deadline.Before(assignment.Deadline) {
			return nil, fmt.Errorf("cannot shorten deadline for published assignment")
		}
	}
	
	// 更新字段
	if req.Title != "" {
		assignment.Title = req.Title
	}
	if req.Description != "" {
		assignment.Description = req.Description
	}
	if !req.Deadline.IsZero() {
		assignment.Deadline = req.Deadline
	}
	if req.TotalScore > 0 {
		assignment.TotalScore = req.TotalScore
	}
	if req.Status != "" {
		assignment.Status = req.Status
		if req.Status == "published" && assignment.PublishedAt == nil {
			now := time.Now()
			assignment.PublishedAt = &now
		}
	}
	
	if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
		return nil, fmt.Errorf("update assignment failed: %w", err)
	}
	
	return assignment, nil
}

// DeleteAssignment 删除作业
func (s *assignmentService) DeleteAssignment(ctx context.Context, id uint, teacherID uint) error {
	// 获取作业并验证权限
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return fmt.Errorf("teacher has no permission to delete this assignment")
	}
	
	// 已有提交的作业不能删除
	stats, err := s.assignmentRepo.GetSubmissionStats(ctx, id)
	if err == nil && stats.SubmittedCount > 0 {
		return fmt.Errorf("cannot delete assignment with submissions")
	}
	
	if err := s.assignmentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete assignment failed: %w", err)
	}
	
	return nil
}

// GetAssignment 获取作业
func (s *assignmentService) GetAssignment(ctx context.Context, id uint) (*model.Assignment, error) {
	return s.assignmentRepo.GetByID(ctx, id)
}

// GetAssignmentDetail 获取作业详情
func (s *assignmentService) GetAssignmentDetail(ctx context.Context, id uint) (*model.AssignmentDetailResponse, error) {
	assignment, err := s.assignmentRepo.GetDetailByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get assignment detail failed: %w", err)
	}
	
	// 转换题目格式
	questions := make([]model.QuestionDetailResponse, len(assignment.Questions))
	for i, q := range assignment.Questions {
		questionDetail := model.QuestionDetailResponse{
			Question: q,
		}
		
		// 处理选择题选项
		if q.Type == model.QuestionTypeChoice && q.Options != "" {
			var options []model.QuestionOption
			if err := json.Unmarshal([]byte(q.Options), &options); err == nil {
				questionDetail.OptionList = options
			}
		}
		
		questions[i] = questionDetail
	}
	
	// 获取统计信息
	stats, err := s.assignmentRepo.GetSubmissionStats(ctx, id)
	if err != nil {
		stats = &model.AssignmentStatistics{} // 使用空统计信息
	}
	
	return &model.AssignmentDetailResponse{
		Assignment:  *assignment,
		Questions:   questions,
		Attachments: assignment.Attachments,
		Statistics:  *stats,
	}, nil
}

// GetTeacherAssignments 获取教师作业列表
func (s *assignmentService) GetTeacherAssignments(ctx context.Context, teacherID uint, page, pageSize int) ([]*model.AssignmentListResponse, int64, error) {
	offset := (page - 1) * pageSize
	assignments, total, err := s.assignmentRepo.GetByTeacherID(ctx, teacherID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("get teacher assignments failed: %w", err)
	}
	
	result := make([]*model.AssignmentListResponse, len(assignments))
	for i, assignment := range assignments {
		// 获取统计信息
		stats, _ := s.assignmentRepo.GetSubmissionStats(ctx, assignment.ID)
		if stats == nil {
			stats = &model.AssignmentStatistics{}
		}
		
		result[i] = &model.AssignmentListResponse{
			ID:             assignment.ID,
			Title:          assignment.Title,
			ClassName:      assignment.Class.ClassName,
			Deadline:       assignment.Deadline,
			TotalScore:     assignment.TotalScore,
			Status:         assignment.Status,
			CreatedAt:      assignment.CreatedAt,
			PublishedAt:    assignment.PublishedAt,
			TotalStudents:  stats.TotalStudents,
			SubmittedCount: stats.SubmittedCount,
			GradedCount:    stats.GradedCount,
			SubmissionRate: stats.SubmissionRate,
		}
	}
	
	return result, total, nil
}

// GetStudentAssignments 获取学生作业列表
func (s *assignmentService) GetStudentAssignments(ctx context.Context, studentID uint, page, pageSize int) ([]*model.StudentAssignmentResponse, int64, error) {
	// TODO: 实现学生作业列表获取逻辑
	// 这里需要根据学生所在班级获取相关作业，并包含提交状态
	return []*model.StudentAssignmentResponse{}, 0, nil
}

// PublishAssignment 发布作业
func (s *assignmentService) PublishAssignment(ctx context.Context, id uint, teacherID uint) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return fmt.Errorf("teacher has no permission to publish this assignment")
	}
	
	if assignment.Status == "published" {
		return fmt.Errorf("assignment is already published")
	}
	
	assignment.Status = "published"
	now := time.Now()
	assignment.PublishedAt = &now
	
	return s.assignmentRepo.Update(ctx, assignment)
}

// UnpublishAssignment 取消发布作业
func (s *assignmentService) UnpublishAssignment(ctx context.Context, id uint, teacherID uint) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return fmt.Errorf("teacher has no permission to unpublish this assignment")
	}
	
	// 检查是否有学生已提交
	stats, err := s.assignmentRepo.GetSubmissionStats(ctx, id)
	if err == nil && stats.SubmittedCount > 0 {
		return fmt.Errorf("cannot unpublish assignment with submissions")
	}
	
	assignment.Status = "draft"
	assignment.PublishedAt = nil
	
	return s.assignmentRepo.Update(ctx, assignment)
}

// GetAssignmentStatistics 获取作业统计信息
func (s *assignmentService) GetAssignmentStatistics(ctx context.Context, id uint, teacherID uint) (*model.AssignmentStatistics, error) {
	// 验证权限
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return nil, fmt.Errorf("teacher has no permission to view statistics for this assignment")
	}
	
	return s.assignmentRepo.GetSubmissionStats(ctx, id)
}