package repository

import (
	"ai-course/internal/model"
	"context"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, roleId string) error
	GetById(ctx context.Context, roleId string) (*model.Role, error)
	GetAll(ctx context.Context) ([]*model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

// Create implements RoleRepository.
func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// Delete implements RoleRepository.
func (r *roleRepository) Delete(ctx context.Context, roleId string) error {
	return r.db.WithContext(ctx).Delete(&model.Role{}, roleId).Error
}

// GetAll implements RoleRepository.
func (r *roleRepository) GetAll(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetById implements RoleRepository.
func (r *roleRepository) GetById(ctx context.Context, roleId string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where("role_id = ?", roleId).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update implements RoleRepository.
func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}
