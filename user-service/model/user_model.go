package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Photo     string    `gorm:"type:varchar(255)" json:"photo"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Roles     []Role    `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
