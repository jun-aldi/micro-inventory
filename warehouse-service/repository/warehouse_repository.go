package repository

import (
	"context"
	"errors"
	"micro-inventory/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseRepositoryInterface interface {
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, warehouseID uint) error
	GetWarehouseByID(ctx context.Context, warehouseID uint) (*model.Warehouse, error)
	GetAllWarehouses(ctx context.Context, page, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error)
}

type WarehouseRepository struct {
	db *gorm.DB
}

// CreateWarehouse implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] CreateWarehouse -1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(warehouse).Error
	}
}

// DeleteWarehouse implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) DeleteWarehouse(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] DeleteWarehouse -1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouse := &model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ", id).Preload("Warehouses").First(modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] DeleteWarehouse -2: %v", err)
			return err
		}

		if len(modelWarehouse.WarehouseProducts) > 0 {
			log.Errorf("[WarehouseRepository] DeleteWarehouse -3: %v", errors.New("warehouse has products"))
			return errors.New("warehouse has products")
		}

		return w.db.WithContext(ctx).Delete(modelWarehouse).Error
	}
}

// GetWarehouseByID implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) GetWarehouseByID(ctx context.Context, id uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetWarehouseByID -1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelWarehouse := &model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetWarehouseByID -2: %v", err)
			return nil, err
		}
		return modelWarehouse, nil
	}
}

// GetWarehouses implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) GetAllWarehouses(ctx context.Context, page, limit int, search string, sortBy string, sortOrder string) ([]model.Warehouse, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetAllWarehouses -1: %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
		if page < 1 {
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

		query := w.db.Model(&model.Warehouse{})

		if search != "" {
			query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+search+"%", "%"+search+"%")
		}

		var warehouses []model.Warehouse
		var total int64
		if err := query.Count(&total).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouses -2: %v", err)
			return nil, 0, err
		}

		if err := query.Order(sortBy + " " + sortOrder).WithContext(ctx).Preload("WarehouseProducts").Offset(offset).Limit(limit).Find(&warehouses).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouses -3: %v", err)
			return nil, 0, err
		}

		return warehouses, total, nil
	}
}

// UpdateWarehouse implements [WarehouseRepositoryInterface].
func (w *WarehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] UpdateWarehouse -1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingWarehouse := &model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", warehouse.ID).Preload("WarehouseProducts").First(existingWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] UpdateWarehouse -2: %v", err)
			return err
		}

		existingWarehouse.Name = warehouse.Name
		existingWarehouse.Address = warehouse.Address
		existingWarehouse.Phone = warehouse.Phone
		existingWarehouse.Photo = warehouse.Photo

		return w.db.WithContext(ctx).Save(existingWarehouse).Error
	}
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepositoryInterface {
	return &WarehouseRepository{db: db}
}
