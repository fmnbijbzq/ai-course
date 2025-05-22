package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"errors"
)

// AddClassRequest 添加班级请求
type AddClassRequest struct {
	Name        string `json:"name" binding:"required"`       // 班级名称
	Description string `json:"description"`                   // 班级描述
	TeacherID   uint   `json:"teacher_id" binding:"required"` // 教师ID
}

// EditClassRequest 编辑班级请求
type EditClassRequest struct {
	ID          uint   `json:"id" binding:"required"`         // 班级ID
	Name        string `json:"name" binding:"required"`       // 班级名称
	Description string `json:"description"`                   // 班级描述
	TeacherID   uint   `json:"teacher_id" binding:"required"` // 教师ID
}

// ListClassRequest 获取班级列表请求
type ListClassRequest struct {
	Page     int `form:"page" binding:"required,min=1"`      // 页码
	PageSize int `form:"page_size" binding:"required,min=1"` // 每页数量
}

// ClassResponse 班级响应
type ClassResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TeacherID   uint   `json:"teacher_id"`
}

// ClassListResponse 班级列表响应
type ClassListResponse struct {
	Total int64           `json:"total"`
	List  []ClassResponse `json:"list"`
}

// ClassService 班级服务接口
type ClassService interface {
	// Add 添加班级
	Add(req *AddClassRequest) (*model.Class, error)
	// Edit 编辑班级
	Edit(req *EditClassRequest) (*model.Class, error)
	// Delete 删除班级
	Delete(id uint) error
	// List 获取班级列表
	List(page, pageSize int) ([]*model.Class, error)
	GetByID(id uint) (*ClassResponse, error)
}

// classService 班级服务实现
type classService struct{}

// NewClassService 创建班级服务实例
func NewClassService() ClassService {
	return &classService{}
}

// Add 添加班级
func (s *classService) Add(req *AddClassRequest) (*model.Class, error) {
	// TODO: 实现添加班级逻辑
	return nil, nil
}

// Edit 编辑班级
func (s *classService) Edit(req *EditClassRequest) (*model.Class, error) {
	// TODO: 实现编辑班级逻辑
	return nil, nil
}

// Delete 删除班级
func (s *classService) Delete(id uint) error {
	// TODO: 实现删除班级逻辑
	return nil
}

// List 获取班级列表
func (s *classService) List(page, pageSize int) ([]*model.Class, error) {
	// TODO: 实现获取班级列表逻辑
	return nil, nil
}

// GetByID 根据ID获取班级
func (s *classService) GetByID(id uint) (*ClassResponse, error) {
	var class model.Class
	if err := repository.DB.First(&class, id).Error; err != nil {
		return nil, errors.New("class not found")
	}

	return &ClassResponse{
		ID:          class.ID,
		Name:        class.Name,
		Description: class.Description,
		TeacherID:   class.TeacherID,
	}, nil
}
