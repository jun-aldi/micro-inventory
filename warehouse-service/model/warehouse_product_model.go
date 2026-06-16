package model

import (
	"time"
)

type WarehouseProduct struct {
	ID          uint  `json:"id" gorm:"primaryKey;autoIncrement"`
	WarehouseID uint  `json:"warehouse_id" gorm:"not null"`
	ProductID   uint  `json:"product_id" gorm:"not null"`
	Stock       int64 `json:"stock" gorm:"not null"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	Warehouse Warehouse `json:"warehouse" gorm:"foreignKey:WarehouseID"`
}
