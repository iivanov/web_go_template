package dto

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	UID   string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersListResponse []UserResponse
