package response

import "micro-inventory/product-service/pkg/pagination"

type ProductResponse struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Barcode   string           `json:"barcode"`
	Thumbnail string           `json:"thumbnail"`
	About     string           `json:"about"`
	Price     float64          `json:"price"`
	IsPopular bool             `json:"is_popular"`
	Category  CategoryResponse `json:"category"`
}

type GetAllProductsResponse struct {
	Products   []ProductResponse             `json:"products"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}
