package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	StudentID string `json:"student_id" binding:"required,min=5"`
	Name      string `json:"name" binding:"required,min=2"`
	Password  string `json:"password" binding:"required,min=6"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint   `json:"id"`
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
}

// UserService 用户服务接口
type UserService interface {
	Register(req *UserRegisterRequest) (*UserResponse, error)
	Login(req *UserLoginRequest) (*UserResponse, error)
}

// userService 用户服务实现
type userService struct{}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{}
}

// Register 用户注册
func (s *userService) Register(req *UserRegisterRequest) (*UserResponse, error) {
	// 检查学号是否已存在
	var existingUser model.User
	result := repository.DB.Where("student_id = ?", req.StudentID).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("student ID already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		StudentID: req.StudentID,
		Name:      req.Name,
		Password:  string(hashedPassword),
	}

	if err := repository.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		StudentID: user.StudentID,
		Name:      user.Name,
	}, nil
}

// Login 用户登录
func (s *userService) Login(req *UserLoginRequest) (*UserResponse, error) {
	var user model.User
	result := repository.DB.Where("student_id = ?", req.StudentID).First(&user)
	if result.Error != nil {
		return nil, errors.New("student ID not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &UserResponse{
		ID:        user.ID,
		StudentID: user.StudentID,
		Name:      user.Name,
	}, nil
}
