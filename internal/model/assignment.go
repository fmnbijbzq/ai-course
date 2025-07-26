package model

import (
	"time"
	"gorm.io/gorm"
)

// Assignment 作业模型
type Assignment struct {
	gorm.Model
	Title             string    `gorm:"type:varchar(200);not null;comment:作业标题" json:"title"`
	Description       string    `gorm:"type:text;comment:作业说明" json:"description"`
	ClassID           uint      `gorm:"not null;comment:班级ID" json:"class_id"`
	TeacherID         uint      `gorm:"not null;comment:教师ID" json:"teacher_id"`
	Deadline          time.Time `gorm:"not null;comment:截止时间" json:"deadline"`
	TotalScore        int       `gorm:"not null;default:100;comment:总分" json:"total_score"`
	Status            string    `gorm:"type:enum('draft','published','closed');default:'draft';comment:状态" json:"status"`
	PublishedAt       *time.Time `gorm:"comment:发布时间" json:"published_at"`
	GradesPublished   bool      `gorm:"default:false;comment:成绩是否已发布" json:"grades_published"`
	GradesPublishedAt *time.Time `gorm:"comment:成绩发布时间" json:"grades_published_at"`

	// 关联关系
	Class       Class        `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Teacher     User         `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Questions   []Question   `gorm:"foreignKey:AssignmentID" json:"questions,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:AssignmentID" json:"attachments,omitempty"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID" json:"submissions,omitempty"`
}

// TableName 指定表名
func (Assignment) TableName() string {
	return "assignments"
}

// IsPublished 检查作业是否已发布
func (a *Assignment) IsPublished() bool {
	return a.Status == "published"
}

// IsOverdue 检查作业是否已过期
func (a *Assignment) IsOverdue() bool {
	return time.Now().After(a.Deadline)
}