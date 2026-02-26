package repository

import (
	"context"
	"micro-inventory/user-service/model"
)

type UserRepositoryInterface interface {

	// =========================
	// USER
	// =========================
	CreateUser(ctx context.Context, user model.User) (*model.User, error)

	GetAllUsers(ctx context.Context,
		page, limit int,
		search, sortBy, sortOrder string,
	) ([]model.User, int64, error)

	GetUserByID(ctx context.Context, id uint) (*model.User, error)

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	UpdateUser(ctx context.Context, user model.User) error

	DeleteUser(ctx context.Context, id uint) error

	GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error)

	// =========================
	// USER ROLE
	// =========================
	AssignUserToRole(ctx context.Context, userID uint, roleID uint) error

	EditAssignUserToRole(ctx context.Context,
		assignRoleID uint,
		userID uint,
		roleID uint,
	) error

	GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error)

	GetAllUserRoles(ctx context.Context,
		page, limit int,
		search, sortBy, sortOrder string,
	) ([]model.UserRole, int64, error)
}

type userRepository struct {
	db *gorm.DB
}
