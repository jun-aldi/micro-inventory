package controller

import (
	"micro-inventory/user-service/controller/request"
	"micro-inventory/user-service/controller/response"
	"micro-inventory/user-service/model"
	"micro-inventory/user-service/pkg/conv"
	"micro-inventory/user-service/pkg/pagination"
	"micro-inventory/user-service/pkg/validator"
	"micro-inventory/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserControllerInterface interface {
	CreateUser(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error

	GetUserByRoleName(c *fiber.Ctx) error

	AssignUserToRole(c *fiber.Ctx) error
	EditAssignUserToRole(c *fiber.Ctx) error
	GetUserRoleByID(c *fiber.Ctx) error
	GetAllUserRoles(c *fiber.Ctx) error
}

type userController struct {
	userUsecase usecase.UserUsecaseInterface
}

// AssignUserToRole implements [UserControllerInterface].
func (u *userController) AssignUserToRole(c *fiber.Ctx) error {
	ctx := c.Context()
	req := request.AssignUserToRoleRequest{}

	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] AssignUserToRole -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] AssignUserToRole -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := u.userUsecase.AssignUserToRole(ctx, req.UserID, req.RoleID); err != nil {
		log.Errorf("[UserController] AssignUserToRole -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User assigned to role successfully",
	})
}

// CreateUser implements [UserControllerInterface].
func (u *userController) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	req := request.CreateUserRequest{}

	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] CreateUser -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] CreateUser -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userModel := model.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	if err := u.userUsecase.CreateUser(ctx, userModel); err != nil {
		log.Errorf("[UserController] CreateUser -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// DeleteUser implements [UserControllerInterface].
func (u *userController) DeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	userId := conv.StringToUint(id)

	if err := u.userUsecase.DeleteUser(ctx, userId); err != nil {
		log.Errorf("[UserController] DeleteUser -1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})

}

// EditAssignUserToRole implements [UserControllerInterface].
func (u *userController) EditAssignUserToRole(c *fiber.Ctx) error {
	ctx := c.Context()
	req := request.AssignUserToRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] EditAssignUserToRole -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] EditAssignUserToRole -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userRoleIdStr := c.Params("userRoleId")
	userRoleId := conv.StringToUint(userRoleIdStr)

	if err := u.userUsecase.EditAssignUserToRole(ctx, userRoleId, req.UserID, req.RoleID); err != nil {
		log.Errorf("[UserController] EditAssignUserToRole -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// GetAllUserRoles implements [UserControllerInterface].
func (u *userController) GetAllUserRoles(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetAllUserRequest

	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[UserController] GetAllUserRoles -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] GetAllUserRoles -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	users, total, err := u.userUsecase.GetAllUserRoles(ctx, int(req.Page), int(req.Limit), req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[UserController] GetAllUserRoles -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserRoleResponse{}

	for _, user := range users {
		roleName := ""

		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		resp = append(resp, response.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Phone:    user.Phone,
			Photo:    user.Photo,
			RoleName: roleName,
		})
	}

	paginationInfo := pagination.CalculatePagination(int(req.Page), int(req.Limit), int(total))

	response := response.GetAllUserResponse{
		Users:      resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users fetched successfully",
		"data":    response,
	})
}

// GetAllUsers implements [UserControllerInterface].
func (u *userController) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetAllUserRequest

	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[UserController] GetAllUser -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] GetAllUser -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	users, total, err := u.userUsecase.GetAllUsers(ctx, int(req.Page), int(req.Limit), req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[UserController] GetAllUser -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserResponse{}

	for _, user := range users {
		roleName := ""

		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		resp = append(resp, response.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Phone:    user.Phone,
			Photo:    user.Photo,
			RoleName: roleName,
		})
	}

	paginationInfo := pagination.CalculatePagination(int(req.Page), int(req.Limit), int(total))

	response := response.GetAllUserResponse{
		Users:      resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users fetched successfully",
		"data":    response,
	})
}

// GetUserByID implements [UserControllerInterface].
func (u *userController) GetUserByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetUserByRoleName implements [UserControllerInterface].
func (u *userController) GetUserByRoleName(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetUserRoleByID implements [UserControllerInterface].
func (u *userController) GetUserRoleByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateUser implements [UserControllerInterface].
func (u *userController) UpdateUser(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewUserController(userUsecase usecase.UserUsecaseInterface) UserControllerInterface {
	return &userController{
		userUsecase: userUsecase,
	}
}
