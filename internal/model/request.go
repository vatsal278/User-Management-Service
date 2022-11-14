package model

import "time"

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}
type SignUpCredentials struct {
	Name                  string    `json:"name" binding:"required"`
	Email                 string    `json:"email" binding:"required"`
	Password              string    `json:"password" binding:"required"`
	RegistrationDate      string    `json:"registration_date" binding:"required"`
	RegistrationTimestamp time.Time `json:"-"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
