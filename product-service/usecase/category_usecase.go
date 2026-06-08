package usecase

import (
	"context"
	"micro-inventory/product-service/model"
)

type CategoryUsecaseInterface interface {
	GetAllCategories(ctx context.Context, page, limit int, search string, sortBy, sortOrder string) ([]model.Category, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	CreateCategory(ctx context.Context, category *model.Category) error
}
