package model

import (
	"gorm.io/gorm"
)

// QuestionType 题目类型枚举
type QuestionType string

const (
	QuestionTypeChoice     QuestionType = "choice"     // 选择题
	QuestionTypeFillBlank  QuestionType = "fill_blank" // 填空题
	QuestionTypeTrueFalse  QuestionType = "true_false" // 判断题
	QuestionTypeEssay      QuestionType = "essay"      // 简答题
)

// Question 题目模型
type Question struct {
	gorm.Model
	AssignmentID uint         `gorm:"not null;comment:作业ID" json:"assignment_id"`
	Type         QuestionType `gorm:"type:enum('choice','fill_blank','true_false','essay');not null;comment:题目类型" json:"type"`
	Content      string       `gorm:"type:text;not null;comment:题目内容" json:"content"`
	Score        int          `gorm:"not null;default:10;comment:分值" json:"score"`
	Order        int          `gorm:"not null;comment:题目顺序" json:"order"`
	Options      string       `gorm:"type:json;comment:选择题选项JSON" json:"options,omitempty"`      // 选择题选项，JSON格式存储
	CorrectAnswer string       `gorm:"type:text;comment:正确答案" json:"correct_answer,omitempty"`    // 客观题的正确答案
	Reference     string       `gorm:"type:text;comment:参考答案" json:"reference,omitempty"`        // 主观题的参考答案
	Explanation   string       `gorm:"type:text;comment:题目解析" json:"explanation,omitempty"`      // 题目解析

	// 关联关系
	Assignment Assignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Answers    []Answer   `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}

// TableName 指定表名
func (Question) TableName() string {
	return "questions"
}

// IsObjective 判断是否为客观题（可自动判分）
func (q *Question) IsObjective() bool {
	return q.Type == QuestionTypeChoice || q.Type == QuestionTypeFillBlank || q.Type == QuestionTypeTrueFalse
}

// IsSubjective 判断是否为主观题（需人工判分）
func (q *Question) IsSubjective() bool {
	return q.Type == QuestionTypeEssay
}

// QuestionOption 选择题选项结构
type QuestionOption struct {
	Key   string `json:"key"`   // 选项标识 A, B, C, D
	Value string `json:"value"` // 选项内容
}

// ChoiceQuestion 选择题专用结构（用于前端交互）
type ChoiceQuestion struct {
	Question
	OptionList    []QuestionOption `json:"option_list"`
	IsMultiple    bool             `json:"is_multiple"`    // 是否多选
	CorrectKeys   []string         `json:"correct_keys"`   // 正确答案的选项标识
}