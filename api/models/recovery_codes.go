package models

import "database/sql"

type RecoveryCode struct {
	Id         int          `json:"id"`
	UserId     int          `json:"userId"`
	IsRedeemed sql.NullBool `json:"isRedeemed"`
	Code       string       `json:"code" validate:"required"`
}
