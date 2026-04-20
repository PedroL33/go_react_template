package models

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"time"

	"github.com/pquerna/otp/totp"
)

type TwoFactorSetupSession struct {
	UserId       int       `json:"userId" db:"user_id"`
	SecretString string    `json:"secretString" validate:"required" db:"secret_string"`
	CreatedAt    time.Time `json:"createdAt" validate:"omitempty" db:"created_timestamp"`
	Expiration   time.Time `json:"expiration" db:"expiration_timestamp"`
}

func (tf *TwoFactorSetupSession) PopulateSecretStringAndReturnBase64QrCode(accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Dashboard",
		AccountName: accountName,
	})

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}

	png.Encode(&buf, img)
	tf.SecretString = key.Secret()

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
