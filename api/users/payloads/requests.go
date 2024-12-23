package payloads

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Complete2faSetupRequest struct {
	OtpCode string `json:"otp_code"`
}
