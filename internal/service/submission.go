package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"context"
	"fmt"
	"time"
)

// SubmissionService 提交服务接口
type SubmissionService interface {
	// 提交管理
	CreateOrUpdateSubmission(ctx context.Context, req *model.SubmissionRequest, studentID uint) (*model.Submission, error)
	SubmitAssignment(ctx context.Context, assignmentID uint, studentID uint) error
	GetSubmission(ctx context.Context, id uint) (*model.Submission, error)
	GetSubmissionDetail(ctx context.Context, id uint) (*model.Submission, error)
	
	// 学生提交列表
	GetStudentSubmissions(ctx context.Context, studentID uint, page, pageSize int) ([]*model.StudentAssignmentResponse, int64, error)
	GetStudentSubmissionByAssignment(ctx context.Context, assignmentID, studentID uint) (*model.StudentAssignmentResponse, error)
	
	// 教师批改列表
	GetSubmissionsForGrading(ctx context.Context, assignmentID uint, page, pageSize int) ([]*model.SubmissionListResponse, int64, error)
	GetGradingDetail(ctx context.Context, submissionID uint, teacherID uint) (*model.GradingDetailResponse, error)
	
	// 自动判分
	AutoGradeSubmission(ctx context.Context, submissionID uint) error
}

// submissionService 提交服务实现
type submissionService struct {
	submissionRepo repository.SubmissionRepository
	answerRepo     repository.AnswerRepository
	assignmentRepo repository.AssignmentRepository
	questionRepo   repository.QuestionRepository
	questionSvc    QuestionService
}

// NewSubmissionService 创建提交服务实例
func NewSubmissionService(
	submissionRepo repository.SubmissionRepository,
	answerRepo repository.AnswerRepository,
	assignmentRepo repository.AssignmentRepository,
	questionRepo repository.QuestionRepository,
	questionSvc QuestionService,
) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
		answerRepo:     answerRepo,
		assignmentRepo: assignmentRepo,
		questionRepo:   questionRepo,
		questionSvc:    questionSvc,
	}
}

// CreateOrUpdateSubmission 创建或更新提交（保存草稿）
func (s *submissionService) CreateOrUpdateSubmission(ctx context.Context, req *model.SubmissionRequest, studentID uint) (*model.Submission, error) {
	// 验证作业是否存在且已发布
	assignment, err := s.assignmentRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.Status != "published" {
		return nil, fmt.Errorf("assignment is not published")
	}
	
	// 检查是否已过截止时间
	if time.Now().After(assignment.Deadline) && req.Status == model.SubmissionStatusSubmitted {
		return nil, fmt.Errorf("assignment deadline has passed")
	}
	
	// 查找是否已有提交记录
	existingSubmission, err := s.submissionRepo.GetByAssignmentAndStudent(ctx, req.AssignmentID, studentID)
	if err != nil {
		// 如果没有找到，创建新的提交记录
		if err.Error() != "record not found" && err.Error() != "get submission by assignment and student failed: record not found" {
			return nil, fmt.Errorf("failed to check existing submission: %w", err)
		}
		existingSubmission = nil
	}
	
	var submission *model.Submission
	
	if existingSubmission == nil {
		// 创建新提交
		submission = &model.Submission{
			AssignmentID: req.AssignmentID,
			StudentID:    studentID,
			Status:       req.Status,
		}
		
		if req.Status == model.SubmissionStatusSubmitted {
			now := time.Now()
			submission.SubmittedAt = &now
		}
		
		if err := s.submissionRepo.Create(ctx, submission); err != nil {
			return nil, fmt.Errorf("create submission failed: %w", err)
		}
	} else {
		// 更新现有提交
		if existingSubmission.Status == model.SubmissionStatusSubmitted && req.Status == model.SubmissionStatusDraft {
			return nil, fmt.Errorf("cannot change submitted assignment back to draft")
		}
		
		submission = existingSubmission
		submission.Status = req.Status
		
		if req.Status == model.SubmissionStatusSubmitted && submission.SubmittedAt == nil {
			now := time.Now()
			submission.SubmittedAt = &now
		}
		
		if err := s.submissionRepo.Update(ctx, submission); err != nil {
			return nil, fmt.Errorf("update submission failed: %w", err)
		}
	}
	
	// 处理答案
	if err := s.processAnswers(ctx, submission.ID, req.Answers); err != nil {
		return nil, fmt.Errorf("process answers failed: %w", err)
	}
	
	// 如果是提交状态，执行自动判分
	if req.Status == model.SubmissionStatusSubmitted {
		if err := s.AutoGradeSubmission(ctx, submission.ID); err != nil {
			// 自动判分失败不影响提交，只记录错误
			// 可以在这里添加日志记录
		}
	}
	
	return submission, nil
}

// processAnswers 处理答案
func (s *submissionService) processAnswers(ctx context.Context, submissionID uint, answerReqs []model.AnswerRequest) error {
	// 获取现有答案
	existingAnswers, err := s.answerRepo.GetBySubmissionID(ctx, submissionID)
	if err != nil && err.Error() != "record not found" {
		return fmt.Errorf("get existing answers failed: %w", err)
	}
	
	// 创建现有答案的映射
	existingAnswerMap := make(map[uint]*model.Answer)
	for _, answer := range existingAnswers {
		existingAnswerMap[answer.QuestionID] = answer
	}
	
	// 处理新答案
	var answersToUpdate []*model.Answer
	var answersToCreate []*model.Answer
	
	for _, answerReq := range answerReqs {
		if existingAnswer, exists := existingAnswerMap[answerReq.QuestionID]; exists {
			// 更新现有答案
			existingAnswer.Content = answerReq.Content
			answersToUpdate = append(answersToUpdate, existingAnswer)
		} else {
			// 创建新答案
			newAnswer := &model.Answer{
				SubmissionID: submissionID,
				QuestionID:   answerReq.QuestionID,
				Content:      answerReq.Content,
			}
			answersToCreate = append(answersToCreate, newAnswer)
		}
	}
	
	// 批量更新
	if len(answersToUpdate) > 0 {
		if err := s.answerRepo.UpdateBatch(ctx, answersToUpdate); err != nil {
			return fmt.Errorf("update answers failed: %w", err)
		}
	}
	
	// 批量创建
	if len(answersToCreate) > 0 {
		if err := s.answerRepo.CreateBatch(ctx, answersToCreate); err != nil {
			return fmt.Errorf("create answers failed: %w", err)
		}
	}
	
	return nil
}

// SubmitAssignment 提交作业
func (s *submissionService) SubmitAssignment(ctx context.Context, assignmentID uint, studentID uint) error {
	// 获取提交记录
	submission, err := s.submissionRepo.GetByAssignmentAndStudent(ctx, assignmentID, studentID)
	if err != nil {
		return fmt.Errorf("submission not found: %w", err)
	}
	
	if submission.Status == model.SubmissionStatusSubmitted {
		return fmt.Errorf("assignment already submitted")
	}
	
	// 验证作业截止时间
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}
	
	if time.Now().After(assignment.Deadline) {
		return fmt.Errorf("assignment deadline has passed")
	}
	
	// 更新提交状态
	submission.Status = model.SubmissionStatusSubmitted
	now := time.Now()
	submission.SubmittedAt = &now
	
	if err := s.submissionRepo.Update(ctx, submission); err != nil {
		return fmt.Errorf("update submission status failed: %w", err)
	}
	
	// 执行自动判分
	if err := s.AutoGradeSubmission(ctx, submission.ID); err != nil {
		// 自动判分失败不影响提交
		return nil
	}
	
	return nil
}

// GetSubmission 获取提交
func (s *submissionService) GetSubmission(ctx context.Context, id uint) (*model.Submission, error) {
	return s.submissionRepo.GetByID(ctx, id)
}

// GetSubmissionDetail 获取提交详情
func (s *submissionService) GetSubmissionDetail(ctx context.Context, id uint) (*model.Submission, error) {
	return s.submissionRepo.GetDetailByID(ctx, id)
}

// GetStudentSubmissions 获取学生提交列表
func (s *submissionService) GetStudentSubmissions(ctx context.Context, studentID uint, page, pageSize int) ([]*model.StudentAssignmentResponse, int64, error) {
	offset := (page - 1) * pageSize
	submissions, total, err := s.submissionRepo.GetByStudentID(ctx, studentID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("get student submissions failed: %w", err)
	}
	
	result := make([]*model.StudentAssignmentResponse, len(submissions))
	for i, submission := range submissions {
		// 获取题目列表
		questions, err := s.questionRepo.GetByAssignmentIDWithOrder(ctx, submission.AssignmentID)
		if err != nil {
			return nil, 0, fmt.Errorf("get questions failed: %w", err)
		}
		
		// 转换题目格式
		questionDetails := make([]model.QuestionDetailResponse, len(questions))
		for j, q := range questions {
			questionDetails[j] = model.QuestionDetailResponse{
				Question: *q,
			}
		}
		
		result[i] = &model.StudentAssignmentResponse{
			Assignment: submission.Assignment,
			Submission: submission,
			Questions:  questionDetails,
		}
	}
	
	return result, total, nil
}

// GetStudentSubmissionByAssignment 获取学生特定作业的提交
func (s *submissionService) GetStudentSubmissionByAssignment(ctx context.Context, assignmentID, studentID uint) (*model.StudentAssignmentResponse, error) {
	// 获取作业详情
	assignment, err := s.assignmentRepo.GetDetailByID(ctx, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	
	// 获取提交记录（如果存在）
	submission, err := s.submissionRepo.GetByAssignmentAndStudent(ctx, assignmentID, studentID)
	if err != nil && err.Error() != "record not found" && err.Error() != "get submission by assignment and student failed: record not found" {
		return nil, fmt.Errorf("get submission failed: %w", err)
	}
	
	// 转换题目格式
	questionDetails := make([]model.QuestionDetailResponse, len(assignment.Questions))
	for i, q := range assignment.Questions {
		questionDetails[i] = model.QuestionDetailResponse{
			Question: q,
		}
	}
	
	return &model.StudentAssignmentResponse{
		Assignment: *assignment,
		Submission: submission,
		Questions:  questionDetails,
		Attachments: assignment.Attachments,
	}, nil
}

// GetSubmissionsForGrading 获取待批改的提交列表
func (s *submissionService) GetSubmissionsForGrading(ctx context.Context, assignmentID uint, page, pageSize int) ([]*model.SubmissionListResponse, int64, error) {
	offset := (page - 1) * pageSize
	submissions, total, err := s.submissionRepo.GetByAssignmentID(ctx, assignmentID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("get submissions for grading failed: %w", err)
	}
	
	result := make([]*model.SubmissionListResponse, len(submissions))
	for i, submission := range submissions {
		result[i] = &model.SubmissionListResponse{
			ID:          submission.ID,
			StudentID:   submission.StudentID,
			StudentName: submission.Student.Name,
			StudentCode: submission.Student.Code,
			Status:      submission.Status,
			Score:       submission.Score,
			SubmittedAt: submission.SubmittedAt,
			GradedAt:    submission.GradedAt,
		}
	}
	
	return result, total, nil
}

// GetGradingDetail 获取批改详情
func (s *submissionService) GetGradingDetail(ctx context.Context, submissionID uint, teacherID uint) (*model.GradingDetailResponse, error) {
	// 获取提交详情
	submission, err := s.submissionRepo.GetDetailByID(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("submission not found: %w", err)
	}
	
	// 验证教师权限
	if submission.Assignment.TeacherID != teacherID {
		return nil, fmt.Errorf("teacher has no permission to grade this submission")
	}
	
	// 构建题目和答案的映射
	questions := make([]model.QuestionWithAnswer, len(submission.Assignment.Questions))
	answerMap := make(map[uint]*model.Answer)
	
	for _, answer := range submission.Answers {
		answerMap[answer.QuestionID] = &answer
	}
	
	for i, question := range submission.Assignment.Questions {
		questions[i] = model.QuestionWithAnswer{
			Question: question,
			Answer:   answerMap[question.ID],
		}
	}
	
	return &model.GradingDetailResponse{
		Submission: *submission,
		Student:    submission.Student,
		Questions:  questions,
	}, nil
}

// AutoGradeSubmission 自动判分
func (s *submissionService) AutoGradeSubmission(ctx context.Context, submissionID uint) error {
	// 获取提交详情
	submission, err := s.submissionRepo.GetDetailByID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("get submission failed: %w", err)
	}
	
	// 如果已经批改过，不重复判分
	if submission.Status == model.SubmissionStatusGraded {
		return nil
	}
	
	totalScore := 0
	answers := submission.Answers
	
	// 逐题判分
	for _, answer := range answers {
		if answer.Question.IsObjective() {
			// 客观题自动判分
			isCorrect, score, err := s.questionSvc.ValidateAnswer(ctx, answer.QuestionID, answer.Content)
			if err != nil {
				continue // 跳过判分失败的题目
			}
			
			answer.IsCorrect = &isCorrect
			answer.Score = score
			now := time.Now()
			answer.GradedAt = &now
			
			totalScore += score
			
			// 更新答案
			if err := s.answerRepo.Update(ctx, &answer); err != nil {
				return fmt.Errorf("update answer failed: %w", err)
			}
		}
		// 主观题需要教师手动批改，跳过
	}
	
	// 更新提交总分
	submission.Score = totalScore
	if err := s.submissionRepo.Update(ctx, submission); err != nil {
		return fmt.Errorf("update submission score failed: %w", err)
	}
	
	return nil
}