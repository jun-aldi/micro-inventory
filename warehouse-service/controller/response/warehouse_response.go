package response

import "micro-inventory/warehouse-service/pkg/pagination"

type WarehouseResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Photo        string `json:"photo"`
	CountProduct int64  `json:"count_product"`
}

type GetAllWarehouseResponse struct {
	Warehouses []WarehouseResponse           `json:"warehouses"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}

type DetailWarehouseResponseStruct struct {
	ID                uint                       `json:"id"`
	Name              string                     `json:"name"`
	Address           string                     `json:"address"`
	Phone             string                     `json:"phone"`
	Photo             string                     `json:"photo"`
	CountProduct      int64                      `json:"count_product"`
	WarehouseProducts []WarehouseProductResponse `json:"warehouse_products"`
}
