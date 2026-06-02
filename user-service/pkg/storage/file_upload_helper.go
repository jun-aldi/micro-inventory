package storage

import (
	"context"
	"fmt"
	"micro-inventory/user-service/configs"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

const (
	MaxImageSize           = 2 * 1024 * 1024
	AllowedImageExtensions = ".jpg, .jpeg, .png, .webp, .svg"
)

type FileUploadHelper struct {
	storage SupabaseInterface
	cfg     configs.Config
}

func NewFileUploadHelper(storage SupabaseInterface, cfg configs.Config) *FileUploadHelper {
	return &FileUploadHelper{
		storage: storage,
		cfg:     cfg,
	}
}

func (h *FileUploadHelper) UploadPhoto(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error) {
	if err := h.ValidateImageFile(file, MaxImageSize); err != nil {
		log.Errorf("failed to validate image file: %v", err)
		return nil, err
	}

	result, err := h.storage.UploadFile(ctx, file, "images")
	if err != nil {
		log.Errorf("failed to upload file to Supabase: %v", err)
		return nil, err
	}

	return result, nil
}

func (h *FileUploadHelper) ValidateImageFile(file *multipart.FileHeader, maxSize int64) error {
	if !ValidateImageSize(file.Size, maxSize) {
		return fmt.Errorf("file size exceeds the allowes size")
	}

	if !validateFileExtension(getFileExtension(file.Filename), AllowedImageExtensions) {
		return fmt.Errorf("file extension is not allowed")
	}

	return nil
}

func ValidateImageSize(size int64, maxSize int64) bool {
	return size <= maxSize
}

func getFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func validateFileExtension(extension string, allowedExtensions string) bool {
	allowed := strings.Split(allowedExtensions, ",")
	for _, ext := range allowed {
		if strings.TrimSpace(ext) == extension {
			return true
		}
	}
	return false
}
