package dto

type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}
