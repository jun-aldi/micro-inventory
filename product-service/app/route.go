package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, container *Container) {

	api := app.Group("/api/v1")
	categories := api.Group("/categories")
	products := api.Group("/products")

	categories.Get("/", container.CategoryController.FindAll)
	categories.Get("/:id", container.CategoryController.FindByID)
	categories.Post("/", container.CategoryController.Create)
	categories.Put("/:id", container.CategoryController.Update)
	categories.Delete("/:id", container.CategoryController.Delete)

	products.Get("/", container.ProductController.FindAll)
	products.Get("/:id", container.ProductController.FindByID)
	products.Get("/barcode /:barcode ", container.ProductController.GetProductByCode)
	products.Post("/", container.ProductController.Create)
	products.Put("/:id", container.ProductController.Update)
	products.Delete("/:id", container.ProductController.Delete)

}
