package service

import (
	"ai-course/internal/model"
	"ai-course/internal/repository"
	"context"
)

type CreateRoleDTO struct {
	RoleId   string `json:"role_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}

type UpdateRoleDTO struct {
	RoleId   string `json:"role_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}

type RoleResponse struct {
	RoleId   string `json:"role_id"`
	RoleName string `json:"role_name"`
}

type RoleService interface {
	Create(ctx context.Context, dto *CreateRoleDTO) error
	Update(ctx context.Context, dto *UpdateRoleDTO) error
	Delete(ctx context.Context, roleId string) error
	GetById(ctx context.Context, roleId string) (*RoleResponse, error)
	GetAll(ctx context.Context) ([]*RoleResponse, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
}

// Create implements RoleService.
func (r *roleService) Create(ctx context.Context, dto *CreateRoleDTO) error {
	role := &model.Role{
		RoleId:   dto.RoleId,
		RoleName: dto.RoleName,
	}
	return r.roleRepo.Create(ctx, role)
}

// Delete implements RoleService.
func (r *roleService) Delete(ctx context.Context, roleId string) error {
	return r.roleRepo.Delete(ctx, roleId)
}

// GetAll implements RoleService.
func (r *roleService) GetAll(ctx context.Context) ([]*RoleResponse, error) {
	roles, err := r.roleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]*RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = r.toRoleResponse(role)
	}
	return responses, nil
}

// GetById implements RoleService.
func (r *roleService) GetById(ctx context.Context, roleId string) (*RoleResponse, error) {
	role, err := r.roleRepo.GetById(ctx, roleId)
	if err != nil {
		return nil, err
	}
	return r.toRoleResponse(role), nil
}

// Update implements RoleService.
func (r *roleService) Update(ctx context.Context, dto *UpdateRoleDTO) error {
	role := &model.Role{
		RoleId:   dto.RoleId,
		RoleName: dto.RoleName,
	}
	return r.roleRepo.Update(ctx, role)
}

func (r *roleService) toRoleResponse(role *model.Role) *RoleResponse {
	return &RoleResponse{
		RoleId:   role.RoleId,
		RoleName: role.RoleName,
	}
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}
