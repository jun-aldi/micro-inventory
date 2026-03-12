package controller

import (
	"micro-inventory/user-service/controller/request"
	"micro-inventory/user-service/controller/response"
	"micro-inventory/user-service/model"
	"micro-inventory/user-service/pkg/conv"
	"micro-inventory/user-service/pkg/validator"
	"micro-inventory/user-service/usecase"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2/log"
)

type RoleControllerInterface interface {
	CreateRole(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
	GetRoleByID(c *fiber.Ctx) error
	GetAllRoles(c *fiber.Ctx) error
}

type roleController struct {
	roleUsecase usecase.RoleUsecaseInterface
}

func NewRoleController(roleUsecase usecase.RoleUsecaseInterface) RoleControllerInterface {
	return &roleController{
		roleUsecase: roleUsecase,
	}
}

func (r *roleController) CreateRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req := request.CreateRoleRequest{}

	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[RoleController] CreateRole -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[RoleController] CreateRole -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Role{
		Name: req.Name,
	}

	if err := r.roleUsecase.CreateRole(ctx, reqModel); err != nil {
		log.Errorf("[RoleController] CreateRole -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role created successfully",
	})

}

func (r *roleController) UpdateRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req := request.CreateRoleRequest{}

	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[RoleController] UpdateRole -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[RoleController] UpdateRole -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Role{
		ID:   conv.StringToUint(c.Params("id")),
		Name: req.Name,
	}

	if err := r.roleUsecase.UpdateRole(ctx, reqModel); err != nil {
		log.Errorf("[RoleController] UpdateRole -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role updated successfully",
	})
}

func (r *roleController) DeleteRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required",
		})
	}

	id := conv.StringToUint(roleID)

	if err := r.roleUsecase.DeleteRole(ctx, id); err != nil {
		log.Errorf("[RoleController] DeleteRole -2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}

func (r *roleController) GetRoleByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	roleID := c.Params("id")
	if roleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required",
		})
	}

	id := conv.StringToUint(roleID)

	role, err := r.roleUsecase.GetRoleByID(ctx, id)
	if err != nil {
		log.Errorf("[RoleController] GetRoleByID -2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role fetched successfully",
		"data":    role,
	})
}

func (r *roleController) GetAllRoles(c *fiber.Ctx) error {
	ctx := c.Context()

	roles, err := r.roleUsecase.GetAllRoles(ctx)
	if err != nil {
		log.Errorf("[RoleController] GetAllRoles -2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.RoleResponse{}
	for _, role := range roles {
		resp = append(resp, response.RoleResponse{
			ID:        role.ID,
			Name:      role.Name,
			CountUser: len(role.Users),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Roles fetched successfully",
		"data":    resp,
	})
}
