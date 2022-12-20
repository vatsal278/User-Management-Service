package model

import "time"

type PingRequest struct {
	Data string `json:"data" validate:"required"`
}
type SignUpCredentials struct {
	Name                  string    `json:"name" validate:"required"`
	Email                 string    `json:"email" validate:"required,email"`
	Password              string    `json:"password" validate:"required,min=8"`
	RegistrationDate      string    `json:"registration_date" validate:"required"`
	RegistrationTimestamp time.Time `json:"-"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Activate struct {
	UserId string `json:"user_id"`
}
