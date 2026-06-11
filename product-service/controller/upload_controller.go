package controller

import (
	"micro-inventory/product-service/controller/response"
	"micro-inventory/product-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadProductImage(ctx *fiber.Ctx) error
	UploadCategoryImage(ctx *fiber.Ctx) error
}

type uploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// UploadCategoryImage implements [UploadControllerInterface].
func (u *uploadController) UploadCategoryImage(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("image")

	if err != nil {
		log.Errorf("failed to get file: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx.Context(), file, "categories")
	if err != nil {
		log.Errorf("failed to upload file: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	response := response.UploadResponse{
		Url:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})

}

// UploadProductImage implements [UploadControllerInterface].
func (u *uploadController) UploadProductImage(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("image")

	if err != nil {
		log.Errorf("failed to get file: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx.Context(), file, "products")
	if err != nil {
		log.Errorf("failed to upload file: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	response := response.UploadResponse{
		Url:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})

}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &uploadController{
		fileUploadHelper: fileUploadHelper,
	}
}
