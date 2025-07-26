package model

import (
	"time"
	"gorm.io/gorm"
)

// SubmissionStatus 提交状态枚举
type SubmissionStatus string

const (
	SubmissionStatusDraft     SubmissionStatus = "draft"     // 草稿
	SubmissionStatusSubmitted SubmissionStatus = "submitted" // 已提交
	SubmissionStatusGraded    SubmissionStatus = "graded"    // 已批改
)

// Submission 学生作业提交模型
type Submission struct {
	gorm.Model
	AssignmentID uint             `gorm:"not null;comment:作业ID" json:"assignment_id"`
	StudentID    uint             `gorm:"not null;comment:学生ID" json:"student_id"`
	Status       SubmissionStatus `gorm:"type:enum('draft','submitted','graded');default:'draft';comment:提交状态" json:"status"`
	Score        int              `gorm:"default:0;comment:总得分" json:"score"`
	SubmittedAt  *time.Time       `gorm:"comment:提交时间" json:"submitted_at"`
	GradedAt     *time.Time       `gorm:"comment:批改时间" json:"graded_at"`
	GradedBy     uint             `gorm:"comment:批改教师ID" json:"graded_by,omitempty"`
	Feedback     string           `gorm:"type:text;comment:教师反馈" json:"feedback,omitempty"`

	// 关联关系
	Assignment Assignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Student    User       `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	GradeTeacher User     `gorm:"foreignKey:GradedBy" json:"grade_teacher,omitempty"`
	Answers    []Answer   `gorm:"foreignKey:SubmissionID" json:"answers,omitempty"`
}

// TableName 指定表名
func (Submission) TableName() string {
	return "submissions"
}

// IsSubmitted 检查是否已提交
func (s *Submission) IsSubmitted() bool {
	return s.Status == SubmissionStatusSubmitted || s.Status == SubmissionStatusGraded
}

// IsGraded 检查是否已批改
func (s *Submission) IsGraded() bool {
	return s.Status == SubmissionStatusGraded
}

// CanBeModified 检查是否可以修改
func (s *Submission) CanBeModified() bool {
	return s.Status == SubmissionStatusDraft
}