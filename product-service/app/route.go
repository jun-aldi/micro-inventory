package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, container *Container) {

	api := app.Group("/api/v1")
	categories := api.Group("/categories")
	products := api.Group("/products")
	uploads := api.Group("/uploads")

	categories.Get("/", container.CategoryController.GetAllCategories)
	categories.Get("/:id", container.CategoryController.GetCategoryByID)
	categories.Post("/", container.CategoryController.CreateCategory)
	categories.Put("/:id", container.CategoryController.UpdateCategory)
	categories.Delete("/:id", container.CategoryController.DeleteCategory)

	products.Get("/", container.ProductController.GetAllProducts)
	products.Get("/barcode/:barcode", container.ProductController.GetProductByBarcode) // move above /:id to avoid conflict
	products.Get("/:id", container.ProductController.GetProductByID)
	products.Post("/", container.ProductController.CreateProduct)
	products.Put("/:id", container.ProductController.UpdateProduct)
	products.Delete("/:id", container.ProductController.DeleteProduct)

	uploads.Post("/category", container.UploadController.UploadCategoryImage)
	uploads.Post("/product", container.UploadController.UploadProductImage)

}
