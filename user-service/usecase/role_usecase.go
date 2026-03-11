package usecase

import (
	"context"

	"micro-inventory/user-service/model"
	"micro-inventory/user-service/repository"
)

// CRUD
type RoleUsecase interface {
	CreateRole(ctx context.Context, role model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type roleUsecase struct {
	roleRepo repository.RoleRepositoryInterface
}

func (r *roleUsecase) CreateRole(ctx context.Context, role model.Role) (*model.Role, error) {
	return r.roleRepo.CreateRole(ctx, role)
}

func (r *roleUsecase) UpdateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.UpdateRole(ctx, role)
}

func (r *roleUsecase) DeleteRole(ctx context.Context, id uint) error {
	return r.roleRepo.DeleteRole(ctx, id)
}

func (r *roleUsecase) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, id)
}

func (r *roleUsecase) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return r.roleRepo.GetAllRoles(ctx)
}

// constructor
func NewRoleUsecase(roleRepo repository.RoleRepositoryInterface) RoleUsecase {
	return &roleUsecase{
		roleRepo: roleRepo,
	}
}
