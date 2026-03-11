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
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) CreateRole(ctx context.Context, role model.Role) (*model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] CreateRole -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
		if err := r.db.WithContext(ctx).Create(&role).Error; err != nil {
			return nil, err
		}
		return &role, nil
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
			log.Errorf("[Role Repository] DeleteRole -3 : failed to delete role")
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
		modelRole := model.Role{}
		if err := r.db.WithContext(ctx).Where("id = ?", role.ID).First(&modelRole).Error; err != nil {
			log.Errorf("[Role Repository] UpdateRole -2 : %v", err)
			return err
		}

		modelRole.Name = role.Name

		return r.db.WithContext(ctx).Save(&modelRole).Error
	}

}

func (r *roleRepository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] GetRoleByID -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelRole := model.Role{}

		if err := r.db.WithContext(ctx).
			Preload("Users").
			First(&modelRole, id).Error; err != nil {

			log.Errorf("[Role Repository] GetRoleByID -2 : %v", err)
			return nil, err
		}

		return &modelRole, nil
	}
}

func (r *roleRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[Role Repository] GetAllRoles -1 : %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelRole := []model.Role{}
		err := r.db.WithContext(ctx).Preload("Users").Find(&modelRole).Error
		if err != nil {
			log.Errorf("[Role Repository] GetAllRoles -2 : %v", err)
			return nil, err
		}
		return modelRole, nil
	}

}
