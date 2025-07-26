package model

import (
	"time"
)

// CreateAssignmentRequest 创建作业请求
type CreateAssignmentRequest struct {
	Title       string             `json:"title" binding:"required,max=200"`
	Description string             `json:"description"`
	ClassID     uint               `json:"class_id" binding:"required"`
	Deadline    time.Time          `json:"deadline" binding:"required"`
	TotalScore  int                `json:"total_score" binding:"min=1"`
	Status      string             `json:"status" binding:"oneof=draft published"`
	Questions   []CreateQuestionRequest `json:"questions"`
}

// UpdateAssignmentRequest 更新作业请求
type UpdateAssignmentRequest struct {
	Title       string    `json:"title" binding:"max=200"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	TotalScore  int       `json:"total_score" binding:"min=1"`
	Status      string    `json:"status" binding:"oneof=draft published closed"`
}

// CreateQuestionRequest 创建题目请求
type CreateQuestionRequest struct {
	Type          QuestionType       `json:"type" binding:"required,oneof=choice fill_blank true_false essay"`
	Content       string             `json:"content" binding:"required"`
	Score         int                `json:"score" binding:"required,min=1"`
	Order         int                `json:"order" binding:"required,min=1"`
	Options       []QuestionOption   `json:"options,omitempty"`
	CorrectAnswer string             `json:"correct_answer,omitempty"`
	Reference     string             `json:"reference,omitempty"`
	Explanation   string             `json:"explanation,omitempty"`
	IsMultiple    bool               `json:"is_multiple,omitempty"` // 选择题是否多选
}

// UpdateQuestionRequest 更新题目请求
type UpdateQuestionRequest struct {
	Content       string           `json:"content"`
	Score         int              `json:"score" binding:"min=1"`
	Order         int              `json:"order" binding:"min=1"`
	Options       []QuestionOption `json:"options,omitempty"`
	CorrectAnswer string           `json:"correct_answer,omitempty"`
	Reference     string           `json:"reference,omitempty"`
	Explanation   string           `json:"explanation,omitempty"`
}

// SubmissionRequest 学生提交答案请求
type SubmissionRequest struct {
	AssignmentID uint             `json:"assignment_id" binding:"required"`
	Answers      []AnswerRequest  `json:"answers" binding:"required"`
	Status       SubmissionStatus `json:"status" binding:"oneof=draft submitted"`
}

// AnswerRequest 答案请求
type AnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// GradeSubmissionRequest 批改作业请求
type GradeSubmissionRequest struct {
	Answers         []GradeAnswerRequest `json:"answers" binding:"required"`
	OverallFeedback string               `json:"overall_feedback"`
}

// GradeAnswerRequest 批改答案请求
type GradeAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Score      int    `json:"score" binding:"min=0"`
	Feedback   string `json:"feedback"`
}

// BatchGradeRequest 批量批改请求
type BatchGradeRequest struct {
	Submissions []BatchGradeItem `json:"submissions" binding:"required"`
}

// BatchGradeItem 批量批改项目
type BatchGradeItem struct {
	SubmissionID    uint                 `json:"submission_id" binding:"required"`
	Answers         []GradeAnswerRequest `json:"answers" binding:"required"`
	OverallFeedback string               `json:"overall_feedback"`
}

// BatchGradeResult 批量批改结果
type BatchGradeResult struct {
	SubmissionID uint   `json:"submission_id"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}

// GradingProgress 批改进度
type GradingProgress struct {
	AssignmentID     uint    `json:"assignment_id"`
	TotalSubmissions int64   `json:"total_submissions"`
	GradedCount      int64   `json:"graded_count"`
	UngradeCount     int64   `json:"ungraded_count"`
	GradingProgress  float64 `json:"grading_progress"`
	GradesPublished  bool    `json:"grades_published"`
}

// SubmissionStatistics 提交统计
type SubmissionStatistics struct {
	TotalSubmissions    int64 `json:"total_submissions"`
	DraftSubmissions    int64 `json:"draft_submissions"`
	SubmittedSubmissions int64 `json:"submitted_submissions"`
	GradedSubmissions   int64 `json:"graded_submissions"`
}

// AssignmentListResponse 作业列表响应
type AssignmentListResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	ClassName   string    `json:"class_name"`
	Deadline    time.Time `json:"deadline"`
	TotalScore  int       `json:"total_score"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	
	// 统计信息
	TotalStudents     int `json:"total_students"`     // 应提交人数
	SubmittedCount    int `json:"submitted_count"`    // 已提交人数
	GradedCount       int `json:"graded_count"`       // 已批改人数
	SubmissionRate    float64 `json:"submission_rate"`    // 提交率
}

// AssignmentDetailResponse 作业详情响应
type AssignmentDetailResponse struct {
	Assignment
	Questions   []QuestionDetailResponse `json:"questions"`
	Attachments []Attachment             `json:"attachments"`
	Statistics  AssignmentStatistics     `json:"statistics"`
}

// QuestionDetailResponse 题目详情响应
type QuestionDetailResponse struct {
	Question
	OptionList  []QuestionOption `json:"option_list,omitempty"`
	IsMultiple  bool             `json:"is_multiple,omitempty"`
	CorrectKeys []string         `json:"correct_keys,omitempty"`
}

// AssignmentStatistics 作业统计信息
type AssignmentStatistics struct {
	TotalStudents     int     `json:"total_students"`
	SubmittedCount    int     `json:"submitted_count"`
	GradedCount       int     `json:"graded_count"`
	SubmissionRate    float64 `json:"submission_rate"`
	AverageScore      float64 `json:"average_score"`
	MaxScore          int     `json:"max_score"`
	MinScore          int     `json:"min_score"`
}

// StudentAssignmentResponse 学生作业响应
type StudentAssignmentResponse struct {
	Assignment  Assignment             `json:"assignment"`
	Submission  *Submission            `json:"submission,omitempty"`
	Questions   []QuestionDetailResponse `json:"questions"`
	Attachments []Attachment           `json:"attachments"`
}

// SubmissionListResponse 提交列表响应
type SubmissionListResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name"`
	StudentCode string    `json:"student_code"`
	Status      SubmissionStatus `json:"status"`
	Score       int       `json:"score"`
	SubmittedAt *time.Time `json:"submitted_at"`
	GradedAt    *time.Time `json:"graded_at"`
}

// SubmissionDetail 提交详情（用于批改）
type SubmissionDetail struct {
	Submission
	Student User     `json:"student"`
	Answers []Answer `json:"answers"`
}

// GradingDetailResponse 批改详情响应
type GradingDetailResponse struct {
	Submission Submission              `json:"submission"`
	Student    User                    `json:"student"`
	Questions  []QuestionWithAnswer    `json:"questions"`
}

// QuestionWithAnswer 题目和答案响应
type QuestionWithAnswer struct {
	Question Question `json:"question"`
	Answer   *Answer  `json:"answer,omitempty"`
}