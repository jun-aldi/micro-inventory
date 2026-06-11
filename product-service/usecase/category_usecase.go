package usecase

import (
	"context"
	"micro-inventory/product-service/model"
	"micro-inventory/product-service/repository"
)

type CategoryUsecaseInterface interface {
	GetAllCategories(ctx context.Context, page, limit int, search string, sortBy, sortOrder string) ([]model.Category, int64, error)
	GetCategoryByID(ctx context.Context, id uint) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	CreateCategory(ctx context.Context, category *model.Category) error
}

type categoryUsecase struct {
	categoryRepository repository.CategoryRepositoryInterface
}

// CreateCategory implements [CategoryUsecaseInterface].
func (c *categoryUsecase) CreateCategory(ctx context.Context, category *model.Category) error {
	return c.categoryRepository.CreateCategory(ctx, category)
}

// DeleteCategory implements [CategoryUsecaseInterface].
func (c *categoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	return c.categoryRepository.DeleteCategory(ctx, id)
}

// GetAllCategories implements [CategoryUsecaseInterface].
func (c *categoryUsecase) GetAllCategories(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Category, int64, error) {
	return c.categoryRepository.GetAllCategories(ctx, page, limit, search, sortBy, sortOrder)
}

// GetCategoryByID implements [CategoryUsecaseInterface].
func (c *categoryUsecase) GetCategoryByID(ctx context.Context, id uint) (*model.Category, error) {
	return c.categoryRepository.GetCategoryByID(ctx, id)
}

// UpdateCategory implements [CategoryUsecaseInterface].
func (c *categoryUsecase) UpdateCategory(ctx context.Context, category *model.Category) error {
	return c.categoryRepository.UpdateCategory(ctx, category)
}

func NewCategoryUsecase(repo repository.CategoryRepositoryInterface) CategoryUsecaseInterface {
	return &categoryUsecase{categoryRepository: repo}
}
