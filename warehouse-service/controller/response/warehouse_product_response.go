package response

type WarehouseProductResponse struct {
	ID          uint              `json:"id"`
	WarehouseID uint              `json:"warehouse_id"`
	ProductID   uint              `json:"product_id"`
	Stock       int               `json:"stock"`
	Warehouse   WarehouseResponse `json:"warehouse"`
}

type GetDetailWarehouseProductByIDResponse struct {
	ID               uint   `json:"id"`
	WarehouseID      uint   `json:"warehouse_id"`
	ProductID        uint   `json:"product_id"`
	Stock            int    `json:"stock"`
	WarehouseName    string `json:"warehouse_name"`
	WarehousePhone   string `json:"warehouse_phone"`
	WarehouseAddress string `json:"warehouse_address"`
	ProductName      string `json:"product_name"`
	ProductBarcode   string `json:"product_barcode"`
	ProductPrice     int    `json:"product_price"`
	ProductThumbnail string `json:"product_thumbnail"`
}

type ProductTotalStockResponse struct {
	ProductID  uint `json:"product_id"`
	TotalStock int  `json:"total_stock"`
}
