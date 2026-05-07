package response

import "micro-inventory/user-service/pkg/pagination"

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	RoleName string `json:"role_name"`
}

type GetAllUserResponse struct {
	Users      []UserResponse                `json:"users"`
	Pagination pagination.PaginationResponse `json:"pagination_data"`
}

type UserRoleResponse struct {
	ID     uint         `json:"id"`
	UserID uint         `json:"user_id"`
	RoleID uint         `json:"role_id"`
	User   UserResponse `json:"user"`
	Role   RoleResponse `json:"role"`
}
