package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"errors"
)

// ClassAddRequest 添加班级请求
type ClassAddRequest struct {
	ClassName string `json:"class_name" binding:"required,min=2,max=50"`
}

// ClassEditRequest 编辑班级请求
type ClassEditRequest struct {
	ClassName string `json:"class_name" binding:"required,min=2,max=50"`
}

// ClassResponse 班级响应
type ClassResponse struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
}

// ClassListResponse 班级列表响应
type ClassListResponse struct {
	Total int64           `json:"total"`
	List  []ClassResponse `json:"list"`
}

// ClassService 班级服务接口
type ClassService interface {
	Add(req *ClassAddRequest) (*ClassResponse, error)
	Edit(id uint, req *ClassEditRequest) (*ClassResponse, error)
	Delete(id uint) error
	List(page, pageSize int) (*ClassListResponse, error)
	GetByID(id uint) (*ClassResponse, error)
}

// classService 班级服务实现
type classService struct{}

// NewClassService 创建班级服务实例
func NewClassService() ClassService {
	return &classService{}
}

// Add 添加班级
func (s *classService) Add(req *ClassAddRequest) (*ClassResponse, error) {
	// 检查班级名是否已存在
	var existingClass model.Class
	result := repository.DB.Where("class_name = ?", req.ClassName).First(&existingClass)
	if result.Error == nil {
		return nil, errors.New("class name already exists")
	}

	// 创建班级
	class := &model.Class{
		ClassName: req.ClassName,
	}

	if err := repository.DB.Create(class).Error; err != nil {
		return nil, err
	}

	return &ClassResponse{
		ID:        class.ID,
		ClassName: class.ClassName,
	}, nil
}

// Edit 编辑班级
func (s *classService) Edit(id uint, req *ClassEditRequest) (*ClassResponse, error) {
	// 检查班级是否存在
	class := &model.Class{}
	if err := repository.DB.First(class, id).Error; err != nil {
		return nil, errors.New("class not found")
	}

	// 检查新名称是否与其他班级重复
	var existingClass model.Class
	result := repository.DB.Where("class_name = ? AND id != ?", req.ClassName, id).First(&existingClass)
	if result.Error == nil {
		return nil, errors.New("class name already exists")
	}

	// 更新班级
	class.ClassName = req.ClassName
	if err := repository.DB.Save(class).Error; err != nil {
		return nil, err
	}

	return &ClassResponse{
		ID:        class.ID,
		ClassName: class.ClassName,
	}, nil
}

// Delete 删除班级
func (s *classService) Delete(id uint) error {
	result := repository.DB.Delete(&model.Class{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("class not found")
	}
	return nil
}

// List 获取班级列表
func (s *classService) List(page, pageSize int) (*ClassListResponse, error) {
	var classes []model.Class
	var total int64

	// 获取总数
	if err := repository.DB.Model(&model.Class{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := repository.DB.Offset(offset).Limit(pageSize).Find(&classes).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	classResponses := make([]ClassResponse, len(classes))
	for i, class := range classes {
		classResponses[i] = ClassResponse{
			ID:        class.ID,
			ClassName: class.ClassName,
		}
	}

	return &ClassListResponse{
		Total: total,
		List:  classResponses,
	}, nil
}

// GetByID 根据ID获取班级
func (s *classService) GetByID(id uint) (*ClassResponse, error) {
	var class model.Class
	if err := repository.DB.First(&class, id).Error; err != nil {
		return nil, errors.New("class not found")
	}

	return &ClassResponse{
		ID:        class.ID,
		ClassName: class.ClassName,
	}, nil
}
