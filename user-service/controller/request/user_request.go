package request

type AssignUserToRoleRequest struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Photo    string `json:"photo" validate:"required"`
}

type GetAllUserRequest struct {
	Page      int64  `json:"page" validate:"omitempty"`
	Limit     int64  `json:"limit" validate:"omitempty"`
	Search    string `json:"search" validate:"omitempty"`
	SortBy    string `json:"sort_by" validate:"omitempty"`
	SortOrder string `json:"sort_order" validate:"omitempty"`
}
