package payloads

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email,lte=60"`
	Password string `json:"password" validate:"required,lte=30,gte=8"`
}

type Complete2faSetupRequest struct {
	OtpCode string `json:"otpCode"`
}

type Disable2faRequest struct {
	Password string `json:"password"`
}

type LoginWithOptCodeRequest struct {
	OtpCode        string `json:"otpCode"`
	LoginSessionId int    `json:"loginSessionId"`
}

type LoginWithRecoveryCodeRequest struct {
	RecoveryCode   string `json:"recoveryCode"`
	LoginSessionId int    `json:"loginSessionId"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword" validate:"required,lte=30,gte=8"`
}
