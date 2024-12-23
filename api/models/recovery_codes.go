package models

import "database/sql"

type RecoveryCode struct {
	Id         int          `json:"id"`
	UserId     int          `json:"user_id"`
	IsRedeemed sql.NullBool `json:"is_redeemed"`
	Code       string       `json:"code" validate:"required"`
}
