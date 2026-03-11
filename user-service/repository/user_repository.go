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

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) AssignUserToRole(ctx context.Context, userID uint, roleID uint) error {

	select {

	case <-ctx.Done():
		log.Errorf("[User Repository] AssignUserToRole -1 : %v", ctx.Err())
		return ctx.Err()
	default:
	}

	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	err := u.db.WithContext(ctx).Create(&userRole).Error
	if err != nil {
		log.Errorf("[User Repository] AssignUserToRole -2 : %v", err)
		return err
	}

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
	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] UpdateUser -1 : %v", ctx.Err())
		return ctx.Err()
	default:
	}

	modelUser := model.User{}

	if err := u.db.WithContext(ctx).Select("id", "name", "email", "password", "photo", "phone").
		Preload("Roles").
		Where("id = ?", user.ID).
		First(&modelUser).Error; err != nil {
		log.Errorf("[User Repository] UpdateUser -2 : %v", err)
		return err
	}

	modelUser.Name = user.Name
	modelUser.Email = user.Email

	if len(user.Password) > 0 {
		modelUser.Password = user.Password
	}

	modelUser.Photo = user.Photo
	modelUser.Phone = user.Phone

	return u.db.WithContext(ctx).Save(&modelUser).Error
}

func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
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

	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] EditAssignUserToRole -1 : %v", ctx.Err())
		return ctx.Err()
	default:
	}

	userRole := model.UserRole{}

	err := u.db.WithContext(ctx).Select("id", "user_id", "role_id").
		Where("id = ?", assignRoleID).
		First(&userRole).Error
	if err != nil {
		log.Errorf("[User Repository] EditAssignUserToRole -2 : %v", err)
		return err
	}

	userRole.UserID = userID
	userRole.RoleID = roleID

	return u.db.WithContext(ctx).Save(&userRole).Error
}

func (u *userRepository) GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[User Repository] GetUserRoleByID -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	userRole := model.UserRole{}

	if err := u.db.WithContext(ctx).Select("id", "user_id", "role_id", "updated_at").
		Preload("User").
		Preload("Role").
		Where("id = ?", assignRoleID).
		First(&userRole).Error; err != nil {
		log.Errorf("[User Repository] GetUserRoleByID -2 : %v", err)
		return nil, err
	}

	return &userRole, nil
}

func (u *userRepository) GetAllUserRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.UserRole, int64, error) {
	// Check if context is done
	select {
	// mengecek apakah context sudah dibatalkan atau timeout
	// ctx time management
	case <-ctx.Done():
		log.Errorf("[User Repository] GetAllUserRoles -1 : %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
	}

	userRoles := []model.UserRole{}

	var totalRecords int64

	// Build Query
	query := u.db.WithContext(ctx).Model(&model.UserRole{})

	// Apply Search Filter if needed
	if search != "" {
		query = query.Joins("JOIN users ON user_roles.user_id = users.id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("users.name ILIKE ? OR users.email ILIKE ? OR roles.name ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get Total Records Count
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Errorf("[User Repository] GetAllUserRoles - Count Error : %v", err)
		return nil, 0, err
	}

	// Apply Sorting
	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "asc"
		}
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("id desc")
	}

	// Apply Pagination
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Execute Query With Preloads
	if err := query.Select("id", "user_id", "role_id").
		Preload("User").
		Preload("Role").
		Find(&userRoles).Error; err != nil {
		log.Errorf("[User Repository] GetAllUserRoles - Find Error : %v", err)
		return nil, 0, err
	}

	return userRoles, totalRecords, nil

}
