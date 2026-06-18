package repository

import (
	"context"
	"micro-inventory/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// Detail warehouse, Get warehouse product by warehouseId and productId,
// Create warehouse product, update warehouse product, delete warehouse product
// delete All products in warehouse product, Get warehouse product by ProductId, Get product total stock,

type WarehouseProductRepositoryInterface interface {
	GetDetailWarehouse(ctx context.Context, warehouseID, productID uint) (*model.Warehouse, error)
	GetDetailWarehouseProductByID(ctx context.Context, warehouseProductID uint) (*model.WarehouseProduct, error)
	GetWarehouseProductsByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error)
	CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	GetWarehouseProductsByWarehouseIDAndProductID(ctx context.Context, warehouseID, productID uint) ([]model.WarehouseProduct, error)
	UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error
	DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error // Fixed: Changed return type to match implementation
	GetProductTotalStock(ctx context.Context, productID uint) (int, error)
}

type WarehouseProductRepository struct {
	db *gorm.DB
}

// CreateWarehouseProduct implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] CreateWarehouseProduct -1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(warehouseProduct).Error
	}
}

// DeleteAllWarehouseProductByProductID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] DeleteAllWarehouseProductbyProductID -1: %v", ctx.Err())
		return ctx.Err()
	default:
		err := w.db.WithContext(ctx).Where("product_id = ? ", productID).Delete(&model.WarehouseProduct{}).Error
		if err != nil {
			log.Errorf("[WarehouseProductRepository] DeleteAllWarehouseProductbyProductID -2: %v", err)
			return err
		}
		return nil
	}
}

// DeleteWarehouseProduct implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] DeleteWarehouseProduct -1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouseProduct := &model.WarehouseProduct{}
		if err := w.db.WithContext(ctx).Where("id = ? ", warehouseProductID).First(modelWarehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] DeleteWarehouseProduct -2: %v", err)
			return err
		}
		return w.db.WithContext(ctx).Delete(modelWarehouseProduct).Error
	}
}

// GetDetailWarehouse implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetDetailWarehouse(ctx context.Context, warehouseID uint, productID uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetDetailWarehouse -1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouse model.Warehouse
		if err := w.db.WithContext(ctx).Where("id = ? ", warehouseID).Select("id", "name", "address", "phone", "photo").Preload("WarehouseProducts").First(&warehouse).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetDetailWarehouse -2: %v", err)
			return nil, err
		}

		return &warehouse, nil
	}
}

// GetProductTotalStock implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetProductTotalStock(ctx context.Context, productID uint) (int, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetProductTotalStock -1: %v", ctx.Err())
		return 0, ctx.Err()
	default:
		var totalStock int64

		if err := w.db.WithContext(ctx).
			Model(&model.WarehouseProduct{}).
			Where("product_id = ?", productID).
			Select("COALESCE(SUM(stock), 0)").
			Scan(&totalStock).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetProductTotalStock -2: %v", err)
			return 0, err
		}

		return int(totalStock), nil
	}
}

// GetDetailWarehouseProductByID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetDetailWarehouseProductByID(ctx context.Context, warehouseProductID uint) (*model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetDetailWarehouseProductByID -1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		warehouseProduct := &model.WarehouseProduct{}
		if err := w.db.WithContext(ctx).Where("id = ? ", warehouseProductID).Select("id", "product_id", "stock", "warehouse_id").Preload("Warehouse").Preload("Product").First(warehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetDetailWarehouseProductByID -2: %v", err)
			return nil, err
		}
		return warehouseProduct, nil
	}
}

// GetWarehouseProductsByProductID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetWarehouseProductsByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetWarehouseProductsByProductID -1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		warehouseProducts := []model.WarehouseProduct{}
		if err := w.db.WithContext(ctx).Where("product_id = ? ", productID).Preload("Warehouse").Preload("Product").Find(&warehouseProducts).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetWarehouseProductsByProductID -2: %v", err)
			return nil, err
		}
		return warehouseProducts, nil
	}
}

// GetWarehouseProductsByWarehouseIDAndProductID implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) GetWarehouseProductsByWarehouseIDAndProductID(ctx context.Context, warehouseID uint, productID uint) ([]model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetWarehouseProductsByWarehouseIDAndProductID -1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		warehouseProducts := []model.WarehouseProduct{}
		if err := w.db.WithContext(ctx).Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).First(&warehouseProducts).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetWarehouseProductsByWarehouseIDAndProductID -2: %v", err)
			return nil, err
		}
		return warehouseProducts, nil
	}
}

// UpdateWarehouseProduct implements [WarehouseProductRepositoryInterface].
func (w *WarehouseProductRepository) UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] UpdateWarehouseProduct -1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingWarehouseProduct := &model.WarehouseProduct{}
		if err := w.db.WithContext(ctx).Where("id = ? ", warehouseProduct.ID).First(existingWarehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] UpdateWarehouseProduct -2: %v", err)
			return err
		}

		existingWarehouseProduct.Stock = warehouseProduct.Stock
		existingWarehouseProduct.WarehouseID = warehouseProduct.WarehouseID
		existingWarehouseProduct.ProductID = warehouseProduct.ProductID

		return w.db.WithContext(ctx).Save(existingWarehouseProduct).Error
	}
}

func NewWarehouseProductRepository(db *gorm.DB) WarehouseProductRepositoryInterface {
	return &WarehouseProductRepository{db: db}
}
