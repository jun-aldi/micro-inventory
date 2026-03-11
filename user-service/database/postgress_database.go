package database

import (
	"fmt"
	"micro-inventory/user-service/configs"
	"micro-inventory/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func ConnectionPostgress(cfg configs.Config) (*Postgres, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.SqlDb.User, cfg.SqlDb.Password, cfg.SqlDb.Host, cfg.SqlDb.Port, cfg.SqlDb.DbName)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Errorf("[Postgres] ConnectionPostgres -1: %v", err)
		return nil, err
	}

	db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[Postgres] ConnectionPostgres -2: %v", err)
		return nil, err
	}

	SeedRole(db)
	SeedManager(db)

	sqlDB.SetMaxOpenConns(cfg.SqlDb.DbMaxOpenCons)
	sqlDB.SetMaxIdleConns(cfg.SqlDb.DbIdleOpenCons)

	return &Postgres{DB: db}, nil
}
