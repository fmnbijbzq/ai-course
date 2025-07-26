package model

import (
	"time"
	"gorm.io/gorm"
)

// Answer 学生答案模型
type Answer struct {
	gorm.Model
	SubmissionID uint   `gorm:"not null;comment:提交记录ID" json:"submission_id"`
	QuestionID   uint   `gorm:"not null;comment:题目ID" json:"question_id"`
	Content      string `gorm:"type:text;comment:答案内容" json:"content"`
	Score        int    `gorm:"default:0;comment:得分" json:"score"`
	IsCorrect    *bool  `gorm:"comment:是否正确(客观题)" json:"is_correct,omitempty"`
	GradedAt     *time.Time `gorm:"comment:批改时间" json:"graded_at,omitempty"`
	Feedback     string `gorm:"type:text;comment:题目反馈" json:"feedback,omitempty"`

	// 关联关系
	Submission Submission `gorm:"foreignKey:SubmissionID" json:"submission,omitempty"`
	Question   Question   `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

// TableName 指定表名
func (Answer) TableName() string {
	return "answers"
}

// IsGraded 检查答案是否已批改
func (a *Answer) IsGraded() bool {
	return a.GradedAt != nil
}