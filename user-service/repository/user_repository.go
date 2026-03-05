package repository

import (
	"context"
	"errors"
	"micro-inventory/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
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

func (u *userRepository) AssignUserToRole(ctx context.Context, userID uint, roleID uint) error {
	return nil
}

func (u *userRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	// Check if context is done
	select {
	// mengecek apakah context sudah dibatalkan atau timeout
	// ctx time management
	case <-ctx.Done():
		log.Errorf("[User Repository] CreateUser -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Errorf("[User Repository] CreateUser -2 : %v", err)
		return nil, err
	}

	// Biar tahu kalau gagal create user
	if user.ID == 0 {
		log.Errorf("[User Repository] CreateUser -3 : %v", err)
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

func (u *userRepository) GetAllUsers(
	ctx context.Context,
	page, limit int,
	search, sortBy, sortOrder string,
) ([]model.User, int64, error) {

	// Check context
	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] GetAllUsers -1 : %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
	}

	// Default pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Default sorting
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	offset := (page - 1) * limit

	// Build base query (FIXED: pakai model.User)
	query := u.db.WithContext(ctx).Model(&model.User{})

	// Search (PostgreSQL ILIKE)
	if search != "" {
		query = query.Where(
			"name ILIKE ? OR email ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	// Count total records
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Errorf("[User Repository] GetAllUsers -2 : %v", err)
		return nil, 0, err
	}

	// Get paginated data
	var users []model.User
	if err := query.
		Select("id", "name", "email", "photo", "phone", "created_at", "updated_at").
		Preload("Roles").
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {

		log.Errorf("[User Repository] GetAllUsers -3 : %v", err)
		return nil, 0, err
	}

	return users, totalRecords, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] GetUserByID -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	modelUsers := model.User{}
	if err := u.db.WithContext(ctx).Select("id", "name", "email", "password", "photo", "phone", "created_at", "updated_at").
		Where("id = ?", id).
		Preload("Roles").
		First(&modelUsers).Error; err != nil {
		log.Errorf("[User Repository] GetUserByID -2 : %v", err)
		return nil, err
	}

	return &modelUsers, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {

	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] GetUserByEmail -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	modelUsers := model.User{}
	if err := u.db.WithContext(ctx).Select("id", "name", "email", "password", "photo", "phone", "created_at", "updated_at").
		Where("email = ?", email).
		Preload("Roles").
		First(&modelUsers).Error; err != nil {
		log.Errorf("[User Repository] GetUserByEmail -2 : %v", err)
		return nil, err
	}

	return &modelUsers, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user model.User) error {
	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
	// Check if context is cancelled or timeout
	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] DeleteUser -1 : %v", ctx.Err())
		return ctx.Err()
	default:
	}

	var modelUser model.User

	err := u.db.WithContext(ctx).
		Select("id", "name", "email", "password", "photo", "phone").
		Preload("Roles").
		Where("id = ?", id).
		First(&modelUser).Error

	if err != nil {
		log.Errorf("[User Repository] DeleteUser -2 : %v", err)
		return err
	}

	return u.db.WithContext(ctx).Delete(&modelUser).Error
}

func (u *userRepository) GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error) {

	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] GetUserByRoleName -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	user := []model.User{}

	// Gunakan Sub Query
	subQuery := u.db.Table("user_roles").
		Select("user_role.user_id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", roleName)

	if err := u.db.WithContext(ctx).
		Where("id IN ?", subQuery).
		Preload("Roles").
		Find(&user).Error; err != nil {
		log.Errorf("[User Repository] GetUserByRoleName -2 : %v", err)
		return nil, err
	}

	return user, nil
}

func (u *userRepository) EditAssignUserToRole(ctx context.Context, assignRoleID uint, userID uint, roleID uint) error {
	return nil
}

func (u *userRepository) GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error) {
	return nil, nil
}

func (u *userRepository) GetAllUserRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.UserRole, int64, error) {
	// Check if context is done
	select {
	// mengecek apakah context sudah dibatalkan atau timeout
	// ctx time management
	case <-ctx.Done():
		log.Errorf("[User Repository] DeleteUser -1 : %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if sortBy == "" {
		sortBy = "created_at"
	}

	if sortOrder == "" {
		sortOrder = "desc"
	}

	// calculate offset
	offset := (page - 1) * limit

	//BUild query
	query := u.db.WithContext(ctx).Model(&model.UserRole{})

	// Add search user if provided
	// Unsesitif ilike
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	//Get Total Count
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Errorf("[User Repository] GetAllUserRoles - Count Error : %v", err)
		return nil, 0, err
	}

	// Get Paginated Data
	var userRoles []model.UserRole
	// select id name mail password photo phone created at updated at
	if err := query.Select("id", "name", "email", "password", "photo", "phone", "created_at", "updated_at").
		Preload("Roles").
		Order(sortBy + " " + sortOrder).
		Offset(offset).Limit(limit).
		Find(&userRoles).Error; err != nil {
		log.Errorf("[User Repository] GetAllUserRoles - Find Error : %v", err)
		return nil, 0, err
	}

	return userRoles, totalRecords, nil

}
