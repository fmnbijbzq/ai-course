package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Code     string `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
	RoleId   string `gorm:"type:varchar(20);not null" json:"role_id"`
	Password string `gorm:"type:varchar(100);not null" json:"-"` // json:"-" 表示不在JSON中显示
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
