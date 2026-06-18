package controller

import (
	"micro-inventory/warehouse-service/controller/request"
	"micro-inventory/warehouse-service/model"
	"micro-inventory/warehouse-service/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseControllerInterface interface {
	CreateWarehouse(c *fiber.Ctx) error
	UpdateWarehouse(c *fiber.Ctx) error
	DeleteWarehouse(c *fiber.Ctx) error
	GetWarehouseByID(c *fiber.Ctx) error
	GetAllWarehouses(c *fiber.Ctx) error
}

type WarehouseController struct {
	warehouseUsecase usecase.WarehouseUsecaseInterface
	validate         *validator.Validate // Added validator instance here
}

// Compile-time interface assertion
var _ WarehouseControllerInterface = (*WarehouseController)(nil)

// CreateWarehouse implements [WarehouseControllerInterface].
func (w *WarehouseController) CreateWarehouse(c *fiber.Ctx) error {
	var req request.CreateWarehouseRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse -1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Request Body",
		})
	}

	// Fixed: Use the struct's validator instance and call .Struct()
	if err := w.validate.Struct(req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse -2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	warehouse := model.Warehouse{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	if err := w.warehouseUsecase.CreateWarehouse(c.UserContext(), &warehouse); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse -3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create warehouse",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Warehouse created successfully",
	})
}

// DeleteWarehouse implements [WarehouseControllerInterface].
func (w *WarehouseController) DeleteWarehouse(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetAllWarehouses implements [WarehouseControllerInterface].
func (w *WarehouseController) GetAllWarehouses(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetWarehouseByID implements [WarehouseControllerInterface].
func (w *WarehouseController) GetWarehouseByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateWarehouse implements [WarehouseControllerInterface].
func (w *WarehouseController) UpdateWarehouse(c *fiber.Ctx) error {
	panic("unimplemented")
}

// NewWarehouseController accepts or initializes the engine dependency
func NewWarehouseController(warehouseUsecase usecase.WarehouseUsecaseInterface) WarehouseControllerInterface {
	return &WarehouseController{
		warehouseUsecase: warehouseUsecase,
		validate:         validator.New(), // Instantiated once here
	}
}
