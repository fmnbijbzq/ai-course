package model

import (
	"gorm.io/gorm"
)

// Class 班级模型
type Class struct {
	gorm.Model
	Code        string `gorm:"type:varchar(50);uniqueIndex;not null;comment:班级代码"`
	Name        string `gorm:"type:varchar(100);not null;comment:班级名称"`
	Description string `gorm:"type:text;comment:班级描述"`
	TeacherID   uint   `gorm:"not null;comment:教师ID"`
}

// TableName 指定表名
func (Class) TableName() string {
	return "classes"
}
