package repository

import (
	"context"
	"micro-inventory/user-service/model"

	"gorm.io/gorm"
)

// ROLE : Create Update Delete Get

type RoleRepositoryInterface interface {
	CreateRole(ctx context.Context, role model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Role, int64, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{db: db}
}
