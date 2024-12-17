package util

import (
	"example/dashboard/api/models"
	"example/dashboard/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func CreateToken(config *config.Config, user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Email": user.Email,
			"ID":    user.Id,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	// tokenString, err := token.SignedString([]byte(config.JwtSecret))
	tokenString, err := token.SignedString(config.JwtSecret)

	if err != nil {
		return "", errors.Wrap(err, "util.CreateToken")
	}

	return tokenString, nil
}
