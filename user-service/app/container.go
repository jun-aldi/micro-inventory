package app

import (
	"micro-inventory/user-service/configs"
	"micro-inventory/user-service/controller"
	"micro-inventory/user-service/database"
	"micro-inventory/user-service/repository"
	"micro-inventory/user-service/service"
	"micro-inventory/user-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	RoleController controller.RoleControllerInterface
	UserController controller.UserControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectionPostgress(*config)

	if err != nil {
		log.Fatalf("Failed to connect to database: ", err)
	}

	rabbitMQService, err := service.NewRabbitMQService(*configs.NewConfig())
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: ", err)
	}

	roleRepo := repository.NewRoleRepository(db.DB)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	roleController := controller.NewRoleController(roleUsecase)

	userRepo := repository.NewUserRepository(db.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, rabbitMQService)
	userController := controller.NewUserController(userUsecase)

	return &Container{
		RoleController: roleController,
		UserController: userController,
	}

}
