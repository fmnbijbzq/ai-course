package service

import (
	"ai-course/internal/logger"
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

// GradingService 批改服务接口
type GradingService interface {
	// GetSubmissionsForGrading 获取待批改的提交列表
	GetSubmissionsForGrading(ctx context.Context, assignmentID, teacherID uint, page, pageSize int, status string) ([]*model.SubmissionDetail, int64, error)
	// GradeSubmission 批改单个提交
	GradeSubmission(ctx context.Context, submissionID uint, req *model.GradeSubmissionRequest, teacherID uint) (*model.SubmissionDetail, error)
	// BatchGrade 批量批改
	BatchGrade(ctx context.Context, req *model.BatchGradeRequest, teacherID uint) ([]*model.BatchGradeResult, error)
	// PublishGrades 发布成绩
	PublishGrades(ctx context.Context, assignmentID, teacherID uint) error
	// GetGradingProgress 获取批改进度
	GetGradingProgress(ctx context.Context, assignmentID, teacherID uint) (*model.GradingProgress, error)
}

// gradingService 批改服务实现
type gradingService struct {
	submissionRepo repository.SubmissionRepository
	answerRepo     repository.AnswerRepository
	assignmentRepo repository.AssignmentRepository
	questionRepo   repository.QuestionRepository
}

// NewGradingService 创建批改服务
func NewGradingService(
	submissionRepo repository.SubmissionRepository,
	answerRepo repository.AnswerRepository,
	assignmentRepo repository.AssignmentRepository,
	questionRepo repository.QuestionRepository,
) GradingService {
	return &gradingService{
		submissionRepo: submissionRepo,
		answerRepo:     answerRepo,
		assignmentRepo: assignmentRepo,
		questionRepo:   questionRepo,
	}
}

// GetSubmissionsForGrading 获取待批改的提交列表
func (s *gradingService) GetSubmissionsForGrading(ctx context.Context, assignmentID, teacherID uint, page, pageSize int, status string) ([]*model.SubmissionDetail, int64, error) {
	// 验证作业是否存在且教师有权限
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		logger.Logger.Error("Failed to get assignment for grading",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return nil, 0, errors.New("assignment not found")
	}

	if assignment.TeacherID != teacherID {
		logger.Logger.Warn("Teacher has no permission to grade assignment",
			zap.Uint("assignment_id", assignmentID),
			zap.Uint("teacher_id", teacherID),
			zap.Uint("assignment_teacher_id", assignment.TeacherID),
		)
		return nil, 0, errors.New("teacher has no permission to grade this assignment")
	}

	// 获取提交列表
	submissions, total, err := s.submissionRepo.GetByAssignmentIDWithPagination(ctx, assignmentID, page, pageSize, status)
	if err != nil {
		logger.Logger.Error("Failed to get submissions for grading",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return nil, 0, err
	}

	return submissions, total, nil
}

// GradeSubmission 批改单个提交
func (s *gradingService) GradeSubmission(ctx context.Context, submissionID uint, req *model.GradeSubmissionRequest, teacherID uint) (*model.SubmissionDetail, error) {
	// 获取提交详情
	submission, err := s.submissionRepo.GetByIDWithDetail(ctx, submissionID)
	if err != nil {
		logger.Logger.Error("Failed to get submission for grading",
			zap.Error(err),
			zap.Uint("submission_id", submissionID),
		)
		return nil, errors.New("submission not found")
	}

	// 验证提交状态
	if submission.Status != model.SubmissionStatusSubmitted {
		logger.Logger.Warn("Submission is not submitted yet",
			zap.Uint("submission_id", submissionID),
			zap.String("status", string(submission.Status)),
		)
		return nil, errors.New("submission is not submitted yet")
	}

	// 验证教师权限
	assignment, err := s.assignmentRepo.GetByID(ctx, submission.AssignmentID)
	if err != nil {
		return nil, errors.New("assignment not found")
	}

	if assignment.TeacherID != teacherID {
		logger.Logger.Warn("Teacher has no permission to grade submission",
			zap.Uint("submission_id", submissionID),
			zap.Uint("teacher_id", teacherID),
			zap.Uint("assignment_teacher_id", assignment.TeacherID),
		)
		return nil, errors.New("teacher has no permission to grade this submission")
	}

	// 批改答案
	totalScore := 0
	for _, gradeAnswer := range req.Answers {
		// 查找对应的答案记录
		var targetAnswer *model.Answer
		for i, answer := range submission.Answers {
			if answer.QuestionID == gradeAnswer.QuestionID {
				targetAnswer = &submission.Answers[i]
				break
			}
		}

		if targetAnswer == nil {
			logger.Logger.Warn("Answer not found for question",
				zap.Uint("question_id", gradeAnswer.QuestionID),
				zap.Uint("submission_id", submissionID),
			)
			continue
		}

		// 更新答案得分和评语
		targetAnswer.Score = gradeAnswer.Score
		targetAnswer.Feedback = gradeAnswer.Feedback
		isCorrect := gradeAnswer.Score > 0
		targetAnswer.IsCorrect = &isCorrect

		err := s.answerRepo.Update(ctx, targetAnswer)
		if err != nil {
			logger.Logger.Error("Failed to update answer score",
				zap.Error(err),
				zap.Uint("answer_id", targetAnswer.ID),
			)
			return nil, err
		}

		totalScore += gradeAnswer.Score
	}

	// 更新提交状态和得分
	submission.Submission.Score = totalScore
	submission.Submission.Status = model.SubmissionStatusGraded
	submission.Submission.GradedAt = &time.Time{}
	*submission.Submission.GradedAt = time.Now()
	submission.Submission.Feedback = req.OverallFeedback

	err = s.submissionRepo.Update(ctx, &submission.Submission)
	if err != nil {
		logger.Logger.Error("Failed to update submission after grading",
			zap.Error(err),
			zap.Uint("submission_id", submissionID),
		)
		return nil, err
	}

	// 重新获取完整的提交详情
	gradedSubmission, err := s.submissionRepo.GetByIDWithDetail(ctx, submissionID)
	if err != nil {
		logger.Logger.Error("Failed to get graded submission detail",
			zap.Error(err),
			zap.Uint("submission_id", submissionID),
		)
		return nil, err
	}

	logger.Logger.Info("Submission graded successfully",
		zap.Uint("submission_id", submissionID),
		zap.Uint("teacher_id", teacherID),
		zap.Int("total_score", totalScore),
	)

	return gradedSubmission, nil
}

// BatchGrade 批量批改
func (s *gradingService) BatchGrade(ctx context.Context, req *model.BatchGradeRequest, teacherID uint) ([]*model.BatchGradeResult, error) {
	results := make([]*model.BatchGradeResult, 0, len(req.Submissions))

	for _, batchItem := range req.Submissions {
		result := &model.BatchGradeResult{
			SubmissionID: batchItem.SubmissionID,
			Success:      false,
		}

		// 构造单个批改请求
		gradeReq := &model.GradeSubmissionRequest{
			Answers:          batchItem.Answers,
			OverallFeedback:  batchItem.OverallFeedback,
		}

		// 执行单个批改
		_, err := s.GradeSubmission(ctx, batchItem.SubmissionID, gradeReq, teacherID)
		if err != nil {
			result.Error = err.Error()
			logger.Logger.Warn("Failed to grade submission in batch",
				zap.Error(err),
				zap.Uint("submission_id", batchItem.SubmissionID),
			)
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	logger.Logger.Info("Batch grading completed",
		zap.Uint("teacher_id", teacherID),
		zap.Int("total_submissions", len(req.Submissions)),
		zap.Int("success_count", successCount),
		zap.Int("failed_count", len(req.Submissions)-successCount),
	)

	return results, nil
}

// PublishGrades 发布成绩
func (s *gradingService) PublishGrades(ctx context.Context, assignmentID, teacherID uint) error {
	// 验证作业是否存在且教师有权限
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		logger.Logger.Error("Failed to get assignment for publishing grades",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return errors.New("assignment not found")
	}

	if assignment.TeacherID != teacherID {
		logger.Logger.Warn("Teacher has no permission to publish grades",
			zap.Uint("assignment_id", assignmentID),
			zap.Uint("teacher_id", teacherID),
			zap.Uint("assignment_teacher_id", assignment.TeacherID),
		)
		return errors.New("teacher has no permission to publish grades for this assignment")
	}

	// 更新作业为已发布成绩状态
	assignment.GradesPublished = true
	assignment.GradesPublishedAt = &time.Time{}
	*assignment.GradesPublishedAt = time.Now()

	err = s.assignmentRepo.Update(ctx, assignment)
	if err != nil {
		logger.Logger.Error("Failed to update assignment grades published status",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return err
	}

	logger.Logger.Info("Grades published successfully",
		zap.Uint("assignment_id", assignmentID),
		zap.Uint("teacher_id", teacherID),
	)

	return nil
}

// GetGradingProgress 获取批改进度
func (s *gradingService) GetGradingProgress(ctx context.Context, assignmentID, teacherID uint) (*model.GradingProgress, error) {
	// 验证作业是否存在且教师有权限
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		logger.Logger.Error("Failed to get assignment for grading progress",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return nil, errors.New("assignment not found")
	}

	if assignment.TeacherID != teacherID {
		logger.Logger.Warn("Teacher has no permission to view grading progress",
			zap.Uint("assignment_id", assignmentID),
			zap.Uint("teacher_id", teacherID),
			zap.Uint("assignment_teacher_id", assignment.TeacherID),
		)
		return nil, errors.New("teacher has no permission to view grading progress for this assignment")
	}

	// 获取提交统计
	stats, err := s.submissionRepo.GetStatistics(ctx, assignmentID)
	if err != nil {
		logger.Logger.Error("Failed to get submission statistics for grading progress",
			zap.Error(err),
			zap.Uint("assignment_id", assignmentID),
		)
		return nil, err
	}

	// 计算批改进度
	progress := &model.GradingProgress{
		AssignmentID:     assignmentID,
		TotalSubmissions: stats.TotalSubmissions,
		GradedCount:      stats.GradedSubmissions,
		UngradeCount:     stats.SubmittedSubmissions,
		GradesPublished:  assignment.GradesPublished,
	}

	if progress.TotalSubmissions > 0 {
		progress.GradingProgress = float64(progress.GradedCount) / float64(progress.TotalSubmissions) * 100
	}

	return progress, nil
}