package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	StudentID string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"student_id"`
	Name      string         `gorm:"type:varchar(50);not null" json:"name"`
	Password  string         `gorm:"type:varchar(100);not null" json:"-"` // json:"-" 表示不在JSON中显示
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
