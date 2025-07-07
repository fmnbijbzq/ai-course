package model

import "time"

type Role struct {
	RoleId   string    `gorm:"type:varchar(20);primaryKey;not null" json:"role_id"`
	RoleName string    `gorm:"type:varchar(50);not null" json:"role_name"`
	CreateAt time.Time `gorm:"type:datetime;not null" json:"create_at"`
	UpdateAt time.Time `gorm:"type:datetime;not null" json:"update_at"`
	DeleteAt time.Time `gorm:"type:datetime;not null" json:"delete_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
