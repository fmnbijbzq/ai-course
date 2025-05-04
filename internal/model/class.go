package model

import (
	"time"

	"gorm.io/gorm"
)

// Class 班级模型
type Class struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ClassName string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"class_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Class) TableName() string {
	return "classes"
}
