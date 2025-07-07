package repository

import (
	"ai-course/internal/model"
	"context"
	"fmt"
	"time"
)

// ClassRepository 班级仓储接口
type ClassRepository interface {
	// Create 创建班级
	Create(ctx context.Context, class *model.Class) error
	// Update 更新班级信息
	Update(ctx context.Context, class *model.Class) error
	// Delete 删除班级
	Delete(ctx context.Context, id uint) error
	// FindByID 根据ID查找班级
	FindByID(ctx context.Context, id uint) (*model.Class, error)
	// FindByCode 根据班级代码查找班级
	FindByCode(ctx context.Context, code string) (*model.Class, error)
	// List 获取班级列表
	List(ctx context.Context, offset, limit int) (int64, []*model.Class, error)
	// GetDB 获取数据库实例
	GetDB() DB
}

// classRepository 班级仓储实现
type classRepository struct {
	db    DB
	cache Cache
}

// NewClassRepository 创建班级仓储实例
func NewClassRepository(db DB, cache Cache) ClassRepository {
	return &classRepository{
		db:    db,
		cache: cache,
	}
}

// Create 创建班级
func (r *classRepository) Create(ctx context.Context, class *model.Class) error {
	return r.db.WithContext(ctx).Create(class)
}

// Update 更新班级信息
func (r *classRepository) Update(ctx context.Context, class *model.Class) error {
	if err := r.db.WithContext(ctx).Save(class); err != nil {
		return err
	}

	// 删除缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("class:id:%d", class.ID)
		r.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// Delete 删除班级
func (r *classRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Class{}, id); err != nil {
		return err
	}

	// 删除缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("class:id:%d", id)
		r.cache.Delete(ctx, cacheKey)
	}

	return nil
}

// FindByID 根据ID查找班级
func (r *classRepository) FindByID(ctx context.Context, id uint) (*model.Class, error) {
	var class model.Class

	// 尝试从缓存获取
	if r.cache != nil {
		cacheKey := fmt.Sprintf("class:id:%d", id)
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if c, ok := cached.(*model.Class); ok {
				return c, nil
			}
		}
	}

	// 从数据库获取
	if err := r.db.WithContext(ctx).First(&class, id); err != nil {
		return nil, err
	}

	// 设置缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("class:id:%d", id)
		r.cache.Set(ctx, cacheKey, &class, time.Hour)
	}

	return &class, nil
}

// FindByCode 根据班级代码查找班级
func (r *classRepository) FindByCode(ctx context.Context, code string) (*model.Class, error) {
	var class model.Class

	// 尝试从缓存获取
	if r.cache != nil {
		cacheKey := fmt.Sprintf("class:code:%s", code)
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if c, ok := cached.(*model.Class); ok {
				return c, nil
			}
		}
	}

	// 从数据库获取
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&class); err != nil {
		return nil, err
	}

	// 设置缓存
	if r.cache != nil {
		cacheKey := fmt.Sprintf("class:code:%s", code)
		r.cache.Set(ctx, cacheKey, &class, time.Hour)
	}

	return &class, nil
}

// List 获取班级列表
func (r *classRepository) List(ctx context.Context, offset, limit int) (int64, []*model.Class, error) {
	var total int64
	var classes []*model.Class

	// 获取总记录数
	if err := r.db.WithContext(ctx).Model(&model.Class{}).Count(&total); err != nil {
		return 0, nil, err
	}

	// 获取分页数据
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&classes); err != nil {
		return 0, nil, err
	}

	return total, classes, nil
}

// GetDB 获取数据库实例
func (r *classRepository) GetDB() DB {
	return r.db
}
