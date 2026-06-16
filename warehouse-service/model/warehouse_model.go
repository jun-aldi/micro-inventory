package model

import (
	"time"

	"gorm.io/gorm"
)

type Warehouse struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Address   string         `json:"address" gorm:"type:text"`
	Photo     string         `json:"photo" gorm:"type:varchar(255)"`
	Phone     string         `json:"phone" gorm:"type:varchar(255)"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	WarehouseProducts []WarehouseProduct `json:"warehouse_products" gorm:"foreignKey:WarehouseID"`
}
