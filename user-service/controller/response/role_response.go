package response

type RoleResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CountUser int    `json:"count_user"`
}
