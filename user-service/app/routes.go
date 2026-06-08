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

	users := api.Group("/users")
	users.Post("/", container.UserController.CreateUser)
	users.Put("/:id", container.UserController.UpdateUser)
	users.Delete("/:id", container.UserController.DeleteUser)
	users.Get("/:id", container.UserController.GetUserByID)
	users.Get("/", container.UserController.GetAllUsers)

	assignedRoles := api.Group("/assign-role")
	assignedRoles.Post("/", container.UserController.AssignUserToRole)
	assignedRoles.Put("/:id", container.UserController.EditAssignUserToRole)
	assignedRoles.Get("/:id", container.UserController.GetUserRoleByID)
	assignedRoles.Get("/", container.UserController.GetAllUserRoles)

	users.Get("/role/:roleName", container.UserController.GetUserByRoleName)

	auth := api.Group("/auth")
	auth.Post("/login", container.AuthController.Login)

	upload := api.Group("/upload")
	upload.Post("/photo", container.UploadController.UploadPhoto)
}
