package model

import (
	"gorm.io/gorm"
)

// Attachment 作业附件模型
type Attachment struct {
	gorm.Model
	AssignmentID uint   `gorm:"not null;comment:作业ID" json:"assignment_id"`
	UploaderID   uint   `gorm:"not null;comment:上传者ID" json:"uploader_id"`
	FileName     string `gorm:"type:varchar(255);not null;comment:文件名" json:"file_name"`
	OriginalName string `gorm:"type:varchar(255);not null;comment:原始文件名" json:"original_name"`
	FilePath     string `gorm:"type:varchar(500);not null;comment:文件路径" json:"file_path"`
	FileSize     int64  `gorm:"not null;comment:文件大小(字节)" json:"file_size"`
	ContentType  string `gorm:"type:varchar(100);not null;comment:文件类型" json:"content_type"`

	// 关联关系
	Assignment Assignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Uploader   User       `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
}

// TableName 指定表名
func (Attachment) TableName() string {
	return "attachments"
}