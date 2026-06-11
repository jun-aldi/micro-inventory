package app

import (
	"log"
	"micro-inventory/product-service/configs"
	"micro-inventory/product-service/controller"
	"micro-inventory/product-service/database"
	"micro-inventory/product-service/pkg/storage"
	"micro-inventory/product-service/repository"
	"micro-inventory/product-service/usecase"
)

type Container struct {
	ProductController  controller.ProductControllerInterface
	CategoryController controller.CategoryControllerInterface
	UploadController   controller.UploadControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectionPostgress(*config)

	if err != nil {
		log.Fatalf("Failed to connect to database: ", err)
	}

	categoryRepo := repository.NewCategoryRepository(db.DB)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryController := controller.NewCategoryController(categoryUsecase)

	productRepo := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productController := controller.NewProductController(productUsecase)

	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewFileUploadHelper(supabaseStorage, *config)
	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		ProductController:  productController,
		CategoryController: categoryController,
		UploadController:   uploadController,
	}

}
