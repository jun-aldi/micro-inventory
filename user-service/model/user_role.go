package model

import (
	"time"
)

type UserRole struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	RoleID    uint `gorm:"not null"`
	User      User `gorm:"foreignKey:UserID"`
	Role      Role `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tabler interface {
	TableName() string
}

func (UserRole) TableName() string {
	return "user_role"
}
