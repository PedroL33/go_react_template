package util

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

func GenerateSecretString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.Wrap(err, "util.ReadRequest")
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
