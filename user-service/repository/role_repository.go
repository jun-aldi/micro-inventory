package repository

import (
	"context"
	"errors"
	"micro-inventory/user-service/model"

	"github.com/gofiber/fiber/v2/log"
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

func (r *roleRepository) CreateRole(ctx context.Context, role model.Role) error {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] CreateRole -1 : %v", ctx.Err())
		return ctx.Err()
	default:
		return r.db.WithContext(ctx).Create(&role).Error
	}

}

func (r *roleRepository) DeleteRole(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] DeleteRole -1 : %v", ctx.Err())
		return ctx.Err()
	default:
		modelRole := model.Role{}

		if err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&modelRole).Error; err != nil {
			log.Errorf("[Role Repository] DeleteRole -2 : %v", err)
			return err
		}

		if len(modelRole.Users) > 0 {
			log.Errorf("[Role Repository] DeleteRole -4 : failed to delete role")
			return errors.New("failed to delete role")
		}

		return r.db.WithContext(ctx).Delete(&modelRole).Error
	}

}

func (r *roleRepository) UpdateRole(ctx context.Context, role model.Role) error {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] UpdateRole -1 : %v", ctx.Err())
		return ctx.Err()
	default:
	}

	modelRole := model.Role{}
	if err := r.db.WithContext(ctx).Where("id = ?", role.ID).First(&modelRole).Error; err != nil {
		log.Errorf("[Role Repository] UpdateRole -2 : %v", err)
		return err
	}

	modelRole.Name = role.Name

	return r.db.WithContext(ctx).Save(&modelRole).Error
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] GetRoleByID -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
	}

	modelRole := model.Role{}
	if err := r.db.WithContext(ctx).Select("id", "name", "created_at", "updated_at").
		Where("id = ?", id).
		First(&modelRole).Error; err != nil {
		log.Errorf("[Role Repository] GetRoleByID -2 : %v", err)
		return nil, err
	}

	return &modelRole, nil
}

func (r *roleRepository) GetAllRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Role, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] GetAllRoles -1 : %v", ctx.Err())
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

	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&model.Role{})

	if search != "" {
		query = query.Where(
			"name ILIKE ?",
			"%"+search+"%",
		)
	}

	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Errorf("[Role Repository] GetAllRoles -2 : %v", err)
		return nil, 0, err
	}

	var roles []model.Role
	if err := query.
		Select("id", "name", "created_at", "updated_at").
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(limit).
		Find(&roles).Error; err != nil {

		log.Errorf("[Role Repository] GetAllRoles -3 : %v", err)
		return nil, 0, err
	}

	return roles, totalRecords, nil
}
