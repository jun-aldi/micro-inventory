package storage

import (
	"context"
	"fmt"
	"micro-inventory/warehouse-service/configs"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error)
}

type SupabaseStorage struct {
	client *storage_go.Client
	cfg    configs.Config
}

// UploadFile implements [SupabaseInterface].
func (s *SupabaseStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate unique file name
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)

	// create file path
	filePath := fmt.Sprintf("%s/%s", folder, fileName)

	// Use the simpler implementation with proper Content-Type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Create clint with proper content-type
	client := storage_go.NewClient(s.cfg.Supabase.Url, s.cfg.Supabase.Key, map[string]string{
		"Content-Type": contentType,
	})

	// Upload file to Supabase storage
	_, err = client.UploadFile(s.cfg.Supabase.Bucket, filePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Supabase: %v", err)
	}

	// Get public URL for the uploaded file
	publicUrl := client.GetPublicUrl(s.cfg.Supabase.Bucket, filePath)

	return &UploadResult{
		URL:      publicUrl.SignedURL,
		Path:     filePath,
		Filename: fileName,
	}, nil
}

type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

func NewSupabaseStorage(cfg configs.Config) SupabaseInterface {
	client := storage_go.NewClient(cfg.Supabase.Url, cfg.Supabase.Key, nil)

	return &SupabaseStorage{
		client: client,
		cfg:    cfg,
	}
}
