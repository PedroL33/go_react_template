package util

import (
	"example/template/rest_server/api/models"
	"example/template/rest_server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func CreateToken(config *config.AppConfig, user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Username":         user.Username,
			"Id":               user.Id,
			"TwoFactorEnabled": user.IsTwoFactorEnabled,
			"RegisteredClaims": jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})

	tokenString, err := token.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return "", Wrap(err)
	}

	return tokenString, nil
}

func VerifyToken(config *config.AppConfig, tokenString string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, Wrap(errors.New("Invalid signing method."))
		}

		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return nil, Wrap(err)
	}

	if !token.Valid {
		return nil, Wrap(errors.New("Invalid token."))
	}

	return claims, nil
}
