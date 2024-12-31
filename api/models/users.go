package models

import (
	"database/sql"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                 int            `json:"id" validate:"omitempty"`
	Email              string         `json:"email" validate:"required,email,lte=60"`
	Password           string         `json:"password" validate:"required,lte=30,gte=8"`
	TwoFactorSecret    sql.NullString `json:"twoFactorSecret" validate:"omitempty"`
	IsTwoFactorEnabled sql.NullBool   `json:"isTwoFactorEnabled" validate:"omitempty"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (u *User) Sanitize() {
	u.Password = ""
	u.TwoFactorSecret = sql.NullString{}
}

func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}
