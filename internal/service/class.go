package service

import (
	"ai-course/internal/model"
	"ai-course/internal/pkg/pagination"
	"ai-course/internal/repository"
	"context"
	"errors"
)

var (
	ErrClassNotFound    = errors.New("班级不存在")
	ErrClassCodeExists  = errors.New("班级代码已存在")
	ErrInvalidClassCode = errors.New("无效的班级代码")
	ErrInvalidClassName = errors.New("无效的班级名称")
)

// CreateClassDTO 创建班级的数据传输对象
type CreateClassDTO struct {
	Code        string `json:"code" binding:"required"`       // 班级代码
	Name        string `json:"name" binding:"required"`       // 班级名称
	Description string `json:"description"`                   // 班级描述
	TeacherID   uint   `json:"teacher_id" binding:"required"` // 教师ID
}

// UpdateClassDTO 更新班级的数据传输对象
type UpdateClassDTO struct {
	ID          uint   `json:"id" binding:"required"`         // 班级ID
	Code        string `json:"code" binding:"required"`       // 班级代码
	Name        string `json:"name" binding:"required"`       // 班级名称
	Description string `json:"description"`                   // 班级描述
	TeacherID   uint   `json:"teacher_id" binding:"required"` // 教师ID
}

// ClassResponse 班级响应
type ClassResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TeacherID   uint   `json:"teacher_id"`
}

// ClassListRequest 获取班级列表的请求参数
type ClassListRequest struct {
	pagination.Params
}

// ClassListResponse 班级列表响应对象
type ClassListResponse struct {
	pagination.Response
}

// ClassService 班级服务接口
type ClassService interface {
	// Create 创建班级
	Create(ctx context.Context, dto *CreateClassDTO) error
	// Update 更新班级信息
	Update(ctx context.Context, dto *UpdateClassDTO) error
	// Delete 删除班级
	Delete(ctx context.Context, id uint) error
	// Get 获取班级信息
	Get(ctx context.Context, id uint) (*ClassResponse, error)
	// GetByCode 根据班级代码获取班级信息
	GetByCode(ctx context.Context, code string) (*ClassResponse, error)
	// List 获取班级列表
	List(ctx context.Context, req *ClassListRequest) (*ClassListResponse, error)
}

// classService 班级服务实现
type classService struct {
	classRepo repository.ClassRepository
}

// NewClassService 创建班级服务实例
func NewClassService(classRepo repository.ClassRepository) ClassService {
	return &classService{
		classRepo: classRepo,
	}
}

// Create 创建班级
func (s *classService) Create(ctx context.Context, dto *CreateClassDTO) error {
	// 验证班级信息
	if err := s.validateCreateDTO(dto); err != nil {
		return err
	}

	// 检查班级代码是否已存在
	existing, err := s.classRepo.FindByCode(ctx, dto.Code)
	if err == nil && existing != nil {
		return ErrClassCodeExists
	}

	// 创建班级实体
	class := &model.Class{
		Code:        dto.Code,
		Name:        dto.Name,
		Description: dto.Description,
		TeacherID:   dto.TeacherID,
	}

	// 创建班级
	return s.classRepo.Create(ctx, class)
}

// Update 更新班级信息
func (s *classService) Update(ctx context.Context, dto *UpdateClassDTO) error {
	// 验证班级信息
	if err := s.validateUpdateDTO(dto); err != nil {
		return err
	}

	// 检查班级是否存在
	existing, err := s.classRepo.FindByID(ctx, dto.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrClassNotFound
	}

	// 如果班级代码发生变化，检查新代码是否已存在
	if existing.Code != dto.Code {
		existingByCode, err := s.classRepo.FindByCode(ctx, dto.Code)
		if err == nil && existingByCode != nil && existingByCode.ID != dto.ID {
			return ErrClassCodeExists
		}
	}

	// 更新班级信息
	existing.Code = dto.Code
	existing.Name = dto.Name
	existing.Description = dto.Description
	existing.TeacherID = dto.TeacherID

	return s.classRepo.Update(ctx, existing)
}

// Delete 删除班级
func (s *classService) Delete(ctx context.Context, id uint) error {
	return s.classRepo.Delete(ctx, id)
}

// Get 获取班级信息
func (s *classService) Get(ctx context.Context, id uint) (*ClassResponse, error) {
	class, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrClassNotFound
	}
	return s.toClassResponse(class), nil
}

// GetByCode 根据班级代码获取班级信息
func (s *classService) GetByCode(ctx context.Context, code string) (*ClassResponse, error) {
	class, err := s.classRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrClassNotFound
	}
	return s.toClassResponse(class), nil
}

// List 获取班级列表
func (s *classService) List(ctx context.Context, req *ClassListRequest) (*ClassListResponse, error) {
	// 验证并设置默认值
	pagination.ValidateAndSetDefaults(&req.Params)

	var classes []*model.Class
	response, err := pagination.Paginate(ctx, s.classRepo.GetDB(), &req.Params, &model.Class{}, &classes)
	if err != nil {
		return nil, err
	}

	// 转换为响应对象
	classResponses := make([]ClassResponse, len(classes))
	for i, class := range classes {
		classResponses[i] = *s.toClassResponse(class)
	}
	response.List = classResponses

	return &ClassListResponse{
		Response: *response,
	}, nil
}

// validateCreateDTO 验证创建班级的数据传输对象
func (s *classService) validateCreateDTO(dto *CreateClassDTO) error {
	if dto == nil {
		return errors.New("班级信息不能为空")
	}
	if dto.Code == "" {
		return ErrInvalidClassCode
	}
	if dto.Name == "" {
		return ErrInvalidClassName
	}
	return nil
}

// validateUpdateDTO 验证更新班级的数据传输对象
func (s *classService) validateUpdateDTO(dto *UpdateClassDTO) error {
	if dto == nil {
		return errors.New("班级信息不能为空")
	}
	if dto.Code == "" {
		return ErrInvalidClassCode
	}
	if dto.Name == "" {
		return ErrInvalidClassName
	}
	return nil
}

// toClassResponse 将班级实体转换为响应对象
func (s *classService) toClassResponse(class *model.Class) *ClassResponse {
	return &ClassResponse{
		ID:          class.ID,
		Code:        class.Code,
		Name:        class.Name,
		Description: class.Description,
		TeacherID:   class.TeacherID,
	}
}
