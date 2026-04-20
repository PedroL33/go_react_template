package models

import "database/sql"

type RecoveryCode struct {
	Id         int          `json:"id" db:"id"`
	UserId     int          `json:"userId" db:"user_id"`
	IsRedeemed sql.NullBool `json:"isRedeemed" db:"is_redeemed"`
	Code       string       `json:"code" validate:"required" db:"code"`
}
