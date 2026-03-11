package database

import (
	"micro-inventory/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) {
	roles := []model.Role{
		{Name: "Manager"},
		{Name: "Keeper"},
	}

	for _, role := range roles {
		if err := db.Create(&role).Error; err != nil {
			log.Errorf("[RoleSeeder] SeedRole -1: %v", err)
		} else {
			log.Infof("[RoleSeeder] SeedRole -1: %v", "Role created succesfully")
		}

	}

}
