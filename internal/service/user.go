package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"ai-course/internal/utils"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUserAlreadyExists  = errors.New("用户已存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrInvalidStudentID   = errors.New("无效的学号")
	ErrInvalidName        = errors.New("无效的用户名")
	ErrInvalidPassword    = errors.New("无效的密码")
)

// CreateUserDTO 创建用户的数据传输对象
type CreateUserDTO struct {
	StudentID string `json:"student_id" binding:"required,min=5"` // 学号
	Name      string `json:"name" binding:"required,min=2"`       // 用户名
	Password  string `json:"password" binding:"required,min=6"`   // 密码
}

// UpdateUserDTO 更新用户的数据传输对象
type UpdateUserDTO struct {
	ID        uint   `json:"id" binding:"required"`               // 用户ID
	StudentID string `json:"student_id" binding:"required,min=5"` // 学号
	Name      string `json:"name" binding:"required,min=2"`       // 用户名
}

// LoginUserDTO 用户登录的数据传输对象
type LoginUserDTO struct {
	StudentID string `json:"student_id" binding:"required"` // 学号
	Password  string `json:"password" binding:"required"`   // 密码
}

// UserResponse 用户响应对象
type UserResponse struct {
	ID        uint   `json:"id"`
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
}

// LoginResponse 登录响应对象
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UserListResponse 用户列表响应对象
type UserListResponse struct {
	Total int64          `json:"total"`
	List  []UserResponse `json:"list"`
}

// UserService 用户服务接口
type UserService interface {
	// Register 用户注册
	Register(ctx context.Context, dto *CreateUserDTO) error
	// Login 用户登录
	Login(ctx context.Context, dto *LoginUserDTO) (*LoginResponse, error)
	// Update 更新用户信息
	Update(ctx context.Context, dto *UpdateUserDTO) error
	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
	// Get 获取用户信息
	Get(ctx context.Context, id uint) (*UserResponse, error)
	// GetByStudentID 根据学号获取用户信息
	GetByStudentID(ctx context.Context, studentID string) (*UserResponse, error)
	// List 获取用户列表
	List(ctx context.Context) (*UserListResponse, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, dto *CreateUserDTO) error {
	// 验证用户信息
	if err := s.validateCreateDTO(dto); err != nil {
		return err
	}

	// 检查用户是否已存在
	existing, err := s.userRepo.FindByStudentID(ctx, dto.StudentID)
	if err == nil && existing != nil {
		return ErrUserAlreadyExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户实体
	user := &model.User{
		Code:     dto.StudentID,
		Name:     dto.Name,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(ctx, user)
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, dto *LoginUserDTO) (*LoginResponse, error) {
	// 根据学号查找用户
	user, err := s.userRepo.FindByStudentID(ctx, dto.StudentID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Code)
	if err != nil {
		return nil, err
	}

	// 构造响应
	return &LoginResponse{
		User:  s.toUserResponse(user),
		Token: token,
	}, nil
}

// Update 更新用户信息
func (s *userService) Update(ctx context.Context, dto *UpdateUserDTO) error {
	// 验证用户信息
	if err := s.validateUpdateDTO(dto); err != nil {
		return err
	}

	// 检查用户是否存在
	existing, err := s.userRepo.FindByID(ctx, dto.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	// 如果学号发生变化，检查新学号是否已存在
	if existing.Code != dto.StudentID {
		existingByStudentID, err := s.userRepo.FindByStudentID(ctx, dto.StudentID)
		if err == nil && existingByStudentID != nil && existingByStudentID.ID != dto.ID {
			return ErrUserAlreadyExists
		}
	}

	// 更新用户信息
	existing.Code = dto.StudentID
	existing.Name = dto.Name

	return s.userRepo.Update(ctx, existing)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// Get 获取用户信息
func (s *userService) Get(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return s.toUserResponse(user), nil
}

// GetByStudentID 根据学号获取用户信息
func (s *userService) GetByStudentID(ctx context.Context, studentID string) (*UserResponse, error) {
	user, err := s.userRepo.FindByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return s.toUserResponse(user), nil
}

// List 获取用户列表
func (s *userService) List(ctx context.Context) (*UserListResponse, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	response := &UserListResponse{
		Total: int64(len(users)),
		List:  make([]UserResponse, len(users)),
	}

	for i, user := range users {
		response.List[i] = *s.toUserResponse(user)
	}

	return response, nil
}

// validateCreateDTO 验证创建用户的数据传输对象
func (s *userService) validateCreateDTO(dto *CreateUserDTO) error {
	if dto == nil {
		return errors.New("用户信息不能为空")
	}
	if dto.StudentID == "" {
		return ErrInvalidStudentID
	}
	if dto.Name == "" {
		return ErrInvalidName
	}
	if dto.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}

// validateUpdateDTO 验证更新用户的数据传输对象
func (s *userService) validateUpdateDTO(dto *UpdateUserDTO) error {
	if dto == nil {
		return errors.New("用户信息不能为空")
	}
	if dto.StudentID == "" {
		return ErrInvalidStudentID
	}
	if dto.Name == "" {
		return ErrInvalidName
	}
	return nil
}

// toUserResponse 将用户实体转换为响应对象
func (s *userService) toUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		StudentID: user.Code,
		Name:      user.Name,
	}
}
