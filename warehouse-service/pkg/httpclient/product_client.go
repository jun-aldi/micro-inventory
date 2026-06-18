package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-inventory/warehouse-service/configs"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
)

type ProductClientInterface interface {
	GetProductByID(ctx context.Context, ProductID uint) (*ProductResponse, error)
	GetProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]ProductResponse, error)
	HealthCheck(ctx context.Context) error
}

type ProductClient struct {
	urlProductService string
	httpClient        *http.Client
}

// GetProductByID implements [ProductClientInterface].
func (p *ProductClient) GetProductByID(ctx context.Context, ProductID uint) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products/%d", p.urlProductService, ProductID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] GetProductByID -1: %v", err)
		return nil, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] GetProductByID -2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[ProductClient] GetProductByID -3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[ProductClient] GetProductByID -4: %s", string(body))
		return nil, errors.New("failed to get product by ID")
	}

	var product ProductResponse
	if err := json.Unmarshal(body, &product); err != nil {
		log.Errorf("[ProductClient] GetProductByID -5: %v", err)
		return nil, err
	}

	return &product, nil
}

// GetProducts implements [ProductClientInterface].
func (p *ProductClient) GetProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products?page=%d&limit=%d&search=%s&sort_by=%s&sort_order=%s", p.urlProductService, page, limit, search, sortBy, sortOrder)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] GetProducts -1: %v", err)
		return nil, err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] GetProducts -2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[ProductClient] GetProducts -3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[ProductClient] GetProducts -4: %s", string(body))
		return nil, errors.New("failed to get products")
	}

	var products ProductListResponse
	if err := json.Unmarshal(body, &products); err != nil {
		log.Errorf("[ProductClient] GetProducts -5: %v", err)
		return nil, err
	}
	return products.Data, nil
}

// HealthCheck implements [ProductClientInterface].
func (p *ProductClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/products/health", p.urlProductService)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] HealthCheck -1: %v", err)
		return err
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] HealthCheck -2: %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to check health")
	}

	return nil
}

type ProductResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	About     string `json:"about"`
	Price     int    `json:"price"`
	Thumbnail string `json:"thumbnail"`
	Category  struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
	} `json:"category" `
}

type ProductListResponse struct {
	Message string            `json:"message"`
	Data    []ProductResponse `json:"data"`
	Error   string            `json:"error,omitempty"`
}

func NewProductClient(httpClient *http.Client, cfg configs.Config) ProductClientInterface {
	return &ProductClient{httpClient: httpClient, urlProductService: cfg.App.UrlProductService}
}
