package app

import (
	"micro-inventory/user-service/configs"
	"micro-inventory/user-service/controller"
	"micro-inventory/user-service/database"
	"micro-inventory/user-service/pkg/storage"
	"micro-inventory/user-service/repository"
	"micro-inventory/user-service/service"
	"micro-inventory/user-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	RoleController   controller.RoleControllerInterface
	UserController   controller.UserControllerInterface
	AuthController   controller.AuthControllerInterface
	UploadController controller.UploadControllerInterface
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

	supabaseStorage := storage.NewSupabaseStorage(*config)

	fileUploadHelper := storage.NewFileUploadHelper(supabaseStorage, *config)

	uploadController := controller.NewUploadController(fileUploadHelper)

	roleRepo := repository.NewRoleRepository(db.DB)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	roleController := controller.NewRoleController(roleUsecase)

	userRepo := repository.NewUserRepository(db.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, rabbitMQService)
	userController := controller.NewUserController(userUsecase)

	authController := controller.NewAuthController(userUsecase)

	return &Container{
		RoleController:   roleController,
		UserController:   userController,
		AuthController:   authController,
		UploadController: uploadController,
	}

}
