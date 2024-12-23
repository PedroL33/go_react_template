package payloads

type Begin2faSetupResponse struct {
	Base64QrCode string `json:"base_64_qrcode"`
}

type Complete2faSetupResponse struct {
	RecoveryCodes []string
}
