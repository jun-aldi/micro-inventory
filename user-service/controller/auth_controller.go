package controller

import (
	"micro-inventory/user-service/controller/request"
	"micro-inventory/user-service/controller/response"
	"micro-inventory/user-service/pkg/conv"
	"micro-inventory/user-service/pkg/validator"
	"micro-inventory/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthControllerInterface interface {
	Login(c *fiber.Ctx) error
}

type authController struct {
	AuthService usecase.UserUsecaseInterface
}

// Login implements AuthControllerInterface.
func (a *authController) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		log.Errorf("[AuthController.Login -1: %v]", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(loginRequest); err != nil {
		log.Errorf("[AuthController.Login -2: %v]", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := a.AuthService.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		log.Errorf("[AuthController.Login -3: %v]", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	if user == nil {
		log.Errorf("[AuthController.Login -4: user not found]")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	isSame := conv.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isSame {
		log.Errorf("[AuthController.Login -5: password not match]")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	loginResponse := response.LoginResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   roles,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successfully",
		"data":    loginResponse,
	})

}

func NewAuthController(authService usecase.UserUsecaseInterface) AuthControllerInterface {
	return &authController{
		AuthService: authService,
	}
}
