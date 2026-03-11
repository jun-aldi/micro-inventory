package database

import (
	"micro-inventory/user-service/model"
	"micro-inventory/user-service/pkg/conv"

	"log"

	"gorm.io/gorm"
)

func SeedManager(db *gorm.DB) {
	bytes, err := conv.HashPassword("manager12345")
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)

	}

	modelRole := model.Role{}

	err = db.Where("name = ?", "Manager").First(&modelRole).Error
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	admin := model.User{
		Name:     "Manager",
		Email:    "manager@mail.com",
		Password: bytes,
		Roles:    []model.Role{modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "manager@mail.com"}).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	} else {
		log.Printf("Admin %s created successfully", admin.Name)
	}

}
