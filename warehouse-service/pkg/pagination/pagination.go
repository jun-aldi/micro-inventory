package pagination

import "math"

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

func CalculatePagination(page, limit, totalRecords int) PaginationResponse {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		CurrentPage:  page,
		Limit:        limit,
		TotalRecords: int64(totalRecords),
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrevious:  page > 1,
	}

}
