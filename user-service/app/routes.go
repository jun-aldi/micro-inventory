package app

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api/v1")
	roles := api.Group("/roles")
	roles.Post("/", container.RoleController.CreateRole)
	roles.Put("/:id", container.RoleController.UpdateRole)
	roles.Delete("/:id", container.RoleController.DeleteRole)
	roles.Get("/:id", container.RoleController.GetRoleByID)
	roles.Get("/", container.RoleController.GetAllRoles)
}
