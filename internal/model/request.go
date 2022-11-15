package model

import "time"

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}
type SignUpCredentials struct {
	Name                  string    `json:"name" validate:"required"`
	Email                 string    `json:"email" validate:"required"`
	Password              string    `json:"password" validate:"required"`
	RegistrationDate      string    `json:"registration_date" validate:"required"`
	RegistrationTimestamp time.Time `json:"-"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
