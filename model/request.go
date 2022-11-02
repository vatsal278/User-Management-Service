package model

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}
type SignUpCredentials struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
