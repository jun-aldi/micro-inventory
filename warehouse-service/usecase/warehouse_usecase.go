package usecase

import (
	"context"
	"micro-inventory/warehouse-service/model"
	"micro-inventory/warehouse-service/repository"
)

type WarehouseUsecaseInterface interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, warehouseID uint) error
	GetWarehouseByID(ctx context.Context, warehouseID uint) (*model.Warehouse, error)
	GetAllWarehouses(ctx context.Context, page, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error)
}

type WarehouseUsecase struct {
	warehouseRepo repository.WarehouseRepositoryInterface
}

// CreateWarehouse implements [WarehouseUsecaseInterface].
func (w *WarehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return w.warehouseRepo.CreateWarehouse(ctx, warehouse)
}

// DeleteWarehouse implements [WarehouseUsecaseInterface].
func (w *WarehouseUsecase) DeleteWarehouse(ctx context.Context, warehouseID uint) error {
	return w.warehouseRepo.DeleteWarehouse(ctx, warehouseID)
}

// GetAllWarehouses implements [WarehouseUsecaseInterface].
func (w *WarehouseUsecase) GetAllWarehouses(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error) {
	return w.warehouseRepo.GetAllWarehouses(ctx, page, limit, search, sortBy, sortOrder)
}

// GetWarehouseByID implements [WarehouseUsecaseInterface].
func (w *WarehouseUsecase) GetWarehouseByID(ctx context.Context, warehouseID uint) (*model.Warehouse, error) {
	return w.warehouseRepo.GetWarehouseByID(ctx, warehouseID)
}

// UpdateWarehouse implements [WarehouseUsecaseInterface].
func (w *WarehouseUsecase) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return w.warehouseRepo.UpdateWarehouse(ctx, warehouse)
}

func NewWarehouseUsecase(warehouseRepo repository.WarehouseRepositoryInterface) WarehouseUsecaseInterface {
	return &WarehouseUsecase{warehouseRepo}
}
