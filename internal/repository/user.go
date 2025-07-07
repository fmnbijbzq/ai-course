package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
	"time"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *model.User) error
	// FindByStudentID 根据学号查找用户
	FindByStudentID(ctx context.Context, studentID string) (*model.User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user *model.User) error
	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*model.User, error)
	// List 获取用户列表
	List(ctx context.Context) ([]*model.User, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db    DB
	cache Cache
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db DB, cache Cache) UserRepository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user)
}

// FindByStudentID 根据学号查找用户
func (r *userRepository) FindByStudentID(ctx context.Context, studentID string) (*model.User, error) {
	var user model.User

	// 尝试从缓存获取
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:student_id:%s", studentID)
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if u, ok := cached.(*model.User); ok {
				return u, nil
			}
		}
	}

	// 从数据库获取
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentID).First(&user); err != nil {
		return nil, err
	}

	// 设置缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:student_id:%s", studentID)
		r.cache.Set(ctx, cacheKey, &user, time.Hour)
	}

	return &user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user); err != nil {
		return err
	}

	// 删除缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:student_id:%s", user.Code)
		r.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	user, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&model.User{}, id); err != nil {
		return err
	}

	// 删除缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:student_id:%s", user.Code)
		r.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	// 尝试从缓存获取
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:id:%d", id)
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if u, ok := cached.(*model.User); ok {
				return u, nil
			}
		}
	}

	// 从数据库获取
	if err := r.db.WithContext(ctx).First(&user, id); err != nil {
		return nil, err
	}

	// 设置缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:id:%d", id)
		r.cache.Set(ctx, cacheKey, &user, time.Hour)
	}

	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}
