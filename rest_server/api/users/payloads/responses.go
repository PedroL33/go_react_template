package payloads

import "example/template/rest_server/api/models"

type Begin2faSetupResponse struct {
	Base64QrCode string `json:"base64Qrcode"`
}

type Complete2faSetupResponse struct {
	RecoveryCodes []string `json:"recoveryCodes"`
}

type LoginResponse struct {
	Token          string       `json:"token,omitempty" validate:"omitempty"`
	User           *models.User `json:"user,omitempty" validate:"omitempty"`
	LoginSessionId int          `json:"loginSessionId,omitempty" validate:"omitempty"`
}

type RegenerateRecoveryCodesResponse struct {
	RecoveryCodes []string `json:"recoveryCodes"`
}
