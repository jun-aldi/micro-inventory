package request

type CreateWarehouseRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}

type GetAllWarehouseRequest struct {
	Page      int    `form:"page" validate:"omitempty,gte=1"`
	Limit     int    `form:"limit" validate:"omitempty,gte=1,lte=100"`
	Search    string `form:"search" validate:"omitempty"`
	SortBy    string `form:"sort_by" validate:"omitempty, oneof=id name address phone photo createdAt"`
	SortOrder string `form:"sort_order" validate:"omitempty, oneof=asc desc"`
}
