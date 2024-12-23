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
			"Email":            user.Email,
			"Id":               user.Id,
			"TwoFactorEnabled": user.IsTwoFactorEnabled,
			"Exp":              time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(config.JwtSecret))
	// tokenString, err := token.SignedString(config.JwtSecret)

	if err != nil {
		return "", errors.Wrap(err, "util.CreateToken")
	}

	return tokenString, nil
}

func VerifyToken(config *config.Config, tokenString string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(errors.New("Invalid signing method."), "util.VerifyToken)")
		}

		return []byte(config.JwtSecret), nil
	})

	if exp, ok := claims["Exp"].(int64); ok {
		if exp < time.Now().Unix() {
			return nil, errors.Wrap(errors.New("Expired token"), "util.VerifyToken")
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "util.VerifyToken")
	}

	if !token.Valid {
		return nil, errors.Wrap(errors.New("Invalid token."), "util.VerifyToken")
	}

	return claims, nil
}
