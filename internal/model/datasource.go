package model

import (
	"time"
)

type DsResponse struct {
	Data string
}

type PingDs struct {
	Data string
}

type User struct {
	Id            string    `json:"id" validate:"required" sql:"id"`
	Email         string    `json:"email" validate:"required" sql:"email"`
	Company       string    `json:"company" validate:"required" sql:"company"`
	Password      string    `json:"password" validate:"required" sql:"password"`
	Name          string    `json:"name" validate:"required" sql:"name"`
	RegisteredOn  time.Time `json:"registered_on" sql:"registered_on"`
	UpdatedOn     time.Time `json:"updated_on" sql:"updated_on"`
	Active        bool      `json:"active" sql:"active"`
	ActiveDevices int       `json:"active_devices" sql:"active_devices"`
	Salt          string    `json:"salt" sql:"salt"`
}

type UserDetails struct {
	Name      string
	Email     string
	Company   string
	LastLogin time.Time
}

const Schema = `
		(
		    user_id varchar(225) not null,
			email varchar(225) not null unique,
		    company_name text,
			name varchar(225) not null,
			password varchar(225) not null,
			registered_on timestamp not null DEFAULT CURRENT_TIMESTAMP,
			updated_on timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		    active boolean not null default false,
		    active_devices int(50) not null default 0,
			salt varchar(225) not null,
			primary key (email)
		);
`
