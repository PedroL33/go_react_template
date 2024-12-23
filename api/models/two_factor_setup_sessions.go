package models

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"time"

	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
)

type TwoFactorSetupSession struct {
	UserId       int       `json:"user_id"`
	SecretString string    `json:"secret_string" validate:"required"`
	CreatedAt    time.Time `json:"created_at" validate:"omitempty"`
	Expiration   time.Time `json:"expiration"`
}

func (tf *TwoFactorSetupSession) PopulateSecretStringAndReturnBase64QrCode(accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Dashboard",
		AccountName: accountName,
	})

	if err != nil {
		return "", errors.Wrap(err, "TwoFactoSetupSession.GenerateBase64QrCode")
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", errors.Wrap(err, "TwoFactoSetupSession.GenerateBase64QrCode")
	}

	png.Encode(&buf, img)
	tf.SecretString = key.Secret()

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
