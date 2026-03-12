package app

import (
	"micro-inventory/user-service/configs"
	"micro-inventory/user-service/controller"
	"micro-inventory/user-service/database"
	"micro-inventory/user-service/repository"
	"micro-inventory/user-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	RoleController controller.RoleControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectionPostgress(*config)

	if err != nil {
		log.Fatalf("Failed to connect to database: ", err)
	}

	roleRepo := repository.NewRoleRepository(db.DB)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	roleController := controller.NewRoleController(roleUsecase)

	return &Container{
		RoleController: roleController,
	}

}
