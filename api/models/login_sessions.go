package models

import (
	"time"
)

type LoginSession struct {
	Id         int       `json:"id"`
	UserId     int       `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	Expiration time.Time `json:"expiration" validate:"required"`
}
