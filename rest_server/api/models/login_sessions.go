package models

import (
	"time"
)

type LoginSession struct {
	Id         int       `json:"id" db:"id"`
	UserId     int       `json:"userId" db:"user_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	Expiration time.Time `json:"expiration" db:"expiration" validate:"required"`
}
