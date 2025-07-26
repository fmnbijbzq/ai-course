package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"context"
	"encoding/json"
	"fmt"
)

// QuestionService 题目服务接口
type QuestionService interface {
	// 题目管理
	CreateQuestion(ctx context.Context, req *model.CreateQuestionRequest, assignmentID uint, teacherID uint) (*model.Question, error)
	UpdateQuestion(ctx context.Context, id uint, req *model.UpdateQuestionRequest, teacherID uint) (*model.Question, error)
	DeleteQuestion(ctx context.Context, id uint, teacherID uint) error
	GetQuestion(ctx context.Context, id uint) (*model.Question, error)
	
	// 题目列表
	GetQuestionsByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.QuestionDetailResponse, error)
	
	// 题目验证
	ValidateAnswer(ctx context.Context, questionID uint, answer string) (bool, int, error)
}

// questionService 题目服务实现
type questionService struct {
	questionRepo   repository.QuestionRepository
	assignmentRepo repository.AssignmentRepository
}

// NewQuestionService 创建题目服务实例
func NewQuestionService(
	questionRepo repository.QuestionRepository,
	assignmentRepo repository.AssignmentRepository,
) QuestionService {
	return &questionService{
		questionRepo:   questionRepo,
		assignmentRepo: assignmentRepo,
	}
}

// CreateQuestion 创建题目
func (s *questionService) CreateQuestion(ctx context.Context, req *model.CreateQuestionRequest, assignmentID uint, teacherID uint) (*model.Question, error) {
	// 验证作业是否存在且教师有权限
	assignment, err := s.assignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return nil, fmt.Errorf("teacher has no permission to add question to this assignment")
	}
	
	// 如果作业已发布，不允许添加题目
	if assignment.Status == "published" {
		return nil, fmt.Errorf("cannot add question to published assignment")
	}
	
	// 创建题目
	question := &model.Question{
		AssignmentID:  assignmentID,
		Type:          req.Type,
		Content:       req.Content,
		Score:         req.Score,
		Order:         req.Order,
		CorrectAnswer: req.CorrectAnswer,
		Reference:     req.Reference,
		Explanation:   req.Explanation,
	}
	
	// 处理选择题选项
	if req.Type == model.QuestionTypeChoice && len(req.Options) > 0 {
		optionsJSON, err := json.Marshal(req.Options)
		if err != nil {
			return nil, fmt.Errorf("marshal question options failed: %w", err)
		}
		question.Options = string(optionsJSON)
		
		// 处理选择题答案
		if req.IsMultiple {
			// 多选题：将正确选项的keys组合成JSON数组
			correctKeys := make([]string, 0)
			for _, option := range req.Options {
				// 这里需要根据实际的选择逻辑来判断哪些选项是正确的
				// 暂时使用简单的逻辑：如果CorrectAnswer包含选项的key，则认为是正确答案
				if contains(req.CorrectAnswer, option.Key) {
					correctKeys = append(correctKeys, option.Key)
				}
			}
			correctKeysJSON, _ := json.Marshal(correctKeys)
			question.CorrectAnswer = string(correctKeysJSON)
		} else {
			// 单选题：直接使用CorrectAnswer作为正确选项的key
			question.CorrectAnswer = req.CorrectAnswer
		}
	}
	
	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, fmt.Errorf("create question failed: %w", err)
	}
	
	return question, nil
}

// UpdateQuestion 更新题目
func (s *questionService) UpdateQuestion(ctx context.Context, id uint, req *model.UpdateQuestionRequest, teacherID uint) (*model.Question, error) {
	// 获取题目并验证权限
	question, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("question not found: %w", err)
	}
	
	// 验证作业权限
	assignment, err := s.assignmentRepo.GetByID(ctx, question.AssignmentID)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return nil, fmt.Errorf("teacher has no permission to update this question")
	}
	
	// 如果作业已发布，限制修改
	if assignment.Status == "published" {
		return nil, fmt.Errorf("cannot update question in published assignment")
	}
	
	// 更新字段
	if req.Content != "" {
		question.Content = req.Content
	}
	if req.Score > 0 {
		question.Score = req.Score
	}
	if req.Order > 0 {
		question.Order = req.Order
	}
	if req.CorrectAnswer != "" {
		question.CorrectAnswer = req.CorrectAnswer
	}
	if req.Reference != "" {
		question.Reference = req.Reference
	}
	if req.Explanation != "" {
		question.Explanation = req.Explanation
	}
	
	// 处理选择题选项更新
	if len(req.Options) > 0 {
		optionsJSON, err := json.Marshal(req.Options)
		if err != nil {
			return nil, fmt.Errorf("marshal question options failed: %w", err)
		}
		question.Options = string(optionsJSON)
	}
	
	if err := s.questionRepo.Update(ctx, question); err != nil {
		return nil, fmt.Errorf("update question failed: %w", err)
	}
	
	return question, nil
}

// DeleteQuestion 删除题目
func (s *questionService) DeleteQuestion(ctx context.Context, id uint, teacherID uint) error {
	// 获取题目并验证权限
	question, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("question not found: %w", err)
	}
	
	// 验证作业权限
	assignment, err := s.assignmentRepo.GetByID(ctx, question.AssignmentID)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}
	
	if assignment.TeacherID != teacherID {
		return fmt.Errorf("teacher has no permission to delete this question")
	}
	
	// 如果作业已发布，不允许删除题目
	if assignment.Status == "published" {
		return fmt.Errorf("cannot delete question from published assignment")
	}
	
	// TODO: 检查是否有学生已经回答了这个题目
	
	if err := s.questionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete question failed: %w", err)
	}
	
	return nil
}

// GetQuestion 获取题目
func (s *questionService) GetQuestion(ctx context.Context, id uint) (*model.Question, error) {
	return s.questionRepo.GetByID(ctx, id)
}

// GetQuestionsByAssignmentID 获取作业的所有题目
func (s *questionService) GetQuestionsByAssignmentID(ctx context.Context, assignmentID uint) ([]*model.QuestionDetailResponse, error) {
	questions, err := s.questionRepo.GetByAssignmentIDWithOrder(ctx, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("get questions by assignment id failed: %w", err)
	}
	
	result := make([]*model.QuestionDetailResponse, len(questions))
	for i, q := range questions {
		questionDetail := &model.QuestionDetailResponse{
			Question: *q,
		}
		
		// 处理选择题选项
		if q.Type == model.QuestionTypeChoice && q.Options != "" {
			var options []model.QuestionOption
			if err := json.Unmarshal([]byte(q.Options), &options); err == nil {
				questionDetail.OptionList = options
				
				// 判断是否为多选题
				if q.CorrectAnswer != "" && q.CorrectAnswer[0] == '[' {
					questionDetail.IsMultiple = true
					var correctKeys []string
					if err := json.Unmarshal([]byte(q.CorrectAnswer), &correctKeys); err == nil {
						questionDetail.CorrectKeys = correctKeys
					}
				} else {
					questionDetail.IsMultiple = false
					questionDetail.CorrectKeys = []string{q.CorrectAnswer}
				}
			}
		}
		
		result[i] = questionDetail
	}
	
	return result, nil
}

// ValidateAnswer 验证答案
func (s *questionService) ValidateAnswer(ctx context.Context, questionID uint, answer string) (bool, int, error) {
	question, err := s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		return false, 0, fmt.Errorf("question not found: %w", err)
	}
	
	// 根据题目类型验证答案
	switch question.Type {
	case model.QuestionTypeChoice:
		return s.validateChoiceAnswer(question, answer)
	case model.QuestionTypeFillBlank:
		return s.validateFillBlankAnswer(question, answer)
	case model.QuestionTypeTrueFalse:
		return s.validateTrueFalseAnswer(question, answer)
	case model.QuestionTypeEssay:
		// 简答题需要人工判分
		return false, 0, nil
	default:
		return false, 0, fmt.Errorf("unsupported question type: %s", question.Type)
	}
}

// validateChoiceAnswer 验证选择题答案
func (s *questionService) validateChoiceAnswer(question *model.Question, answer string) (bool, int, error) {
	if question.CorrectAnswer == "" {
		return false, 0, fmt.Errorf("question has no correct answer")
	}
	
	// 检查是否为多选题
	if question.CorrectAnswer[0] == '[' {
		// 多选题
		var correctKeys []string
		if err := json.Unmarshal([]byte(question.CorrectAnswer), &correctKeys); err != nil {
			return false, 0, fmt.Errorf("parse correct answer failed: %w", err)
		}
		
		var studentKeys []string
		if err := json.Unmarshal([]byte(answer), &studentKeys); err != nil {
			return false, 0, fmt.Errorf("parse student answer failed: %w", err)
		}
		
		// 比较答案
		if len(correctKeys) != len(studentKeys) {
			return false, 0, nil
		}
		
		correctMap := make(map[string]bool)
		for _, key := range correctKeys {
			correctMap[key] = true
		}
		
		for _, key := range studentKeys {
			if !correctMap[key] {
				return false, 0, nil
			}
		}
		
		return true, question.Score, nil
	} else {
		// 单选题
		if answer == question.CorrectAnswer {
			return true, question.Score, nil
		}
		return false, 0, nil
	}
}

// validateFillBlankAnswer 验证填空题答案
func (s *questionService) validateFillBlankAnswer(question *model.Question, answer string) (bool, int, error) {
	if question.CorrectAnswer == "" {
		return false, 0, fmt.Errorf("question has no correct answer")
	}
	
	// 简单的字符串匹配（可以后续增强为支持多个正确答案、忽略大小写等）
	if answer == question.CorrectAnswer {
		return true, question.Score, nil
	}
	
	return false, 0, nil
}

// validateTrueFalseAnswer 验证判断题答案
func (s *questionService) validateTrueFalseAnswer(question *model.Question, answer string) (bool, int, error) {
	if question.CorrectAnswer == "" {
		return false, 0, fmt.Errorf("question has no correct answer")
	}
	
	// 标准化答案格式
	normalizedCorrect := normalizeBoolean(question.CorrectAnswer)
	normalizedAnswer := normalizeBoolean(answer)
	
	if normalizedAnswer == normalizedCorrect {
		return true, question.Score, nil
	}
	
	return false, 0, nil
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)))
}

// 辅助函数：标准化布尔值表示
func normalizeBoolean(s string) string {
	switch s {
	case "true", "True", "TRUE", "1", "对", "正确", "是":
		return "true"
	case "false", "False", "FALSE", "0", "错", "错误", "否":
		return "false"
	default:
		return s
	}
}