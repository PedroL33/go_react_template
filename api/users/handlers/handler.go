package handlers

import (
	"example/dashboard/api/base"
	"example/dashboard/api/middleware"
	"example/dashboard/api/models"
	"example/dashboard/api/users"
	"example/dashboard/api/users/payloads"
	"example/dashboard/config"
	http_errors "example/dashboard/errors"
	"example/dashboard/util"
	"net/http"

	"github.com/gorilla/mux"
)

type usersHandlers struct {
	controller users.Controller
	cfg        *config.Config
	logger     util.Logger
}

func NewUsersHandlers(cfg *config.Config, controller users.Controller, logger util.Logger) users.Handler {
	return &usersHandlers{cfg: cfg, controller: controller, logger: logger}
}

func MapUsersRoutes(h users.Handler, router *mux.Router, m middleware.MiddleWareManager) {
	router.HandleFunc("/user", h.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", h.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/login_otp", h.VerifyLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/begin2fa", m.Auth(h.Begin2faSetup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/complete2fa", m.Auth(h.Complete2faSetup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/disable2fa", m.Auth(h.Disable2fa)).Methods("POST", "OPTIONS")
	router.HandleFunc("/login_recovery_code", h.VerifyLoginWithRecoveryCode).Methods("POST", "OPTIONS")
	router.HandleFunc("/change_password", m.Auth(h.ChangePassword)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/regenerate_recovery_codes", m.Auth(h.RegenerateRecoveryCodes)).Methods("POST", "OPTIONS")
}

func (h *usersHandlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	createUserRequest := &payloads.CreateUserRequest{}
	err := util.ReadRequest(ctx, r, createUserRequest)

	if err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	createdUser, err := h.controller.CreateUser(ctx, createUserRequest)

	if err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, createdUser)
}

func (h *usersHandlers) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var err error
	loginRequest := &payloads.UserLoginRequest{}
	if err = util.ReadRequest(ctx, r, loginRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	var loginResponse *payloads.LoginResponse
	if loginResponse, err = h.controller.Login(ctx, &models.User{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, loginResponse)
}

func (h *usersHandlers) VerifyLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	verifyLoginRequest := &payloads.LoginWithOptCodeRequest{}
	if err = util.ReadRequest(ctx, r, verifyLoginRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	var loginResponse *payloads.LoginResponse
	if loginResponse, err = h.controller.VerifyLogin(ctx, verifyLoginRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, loginResponse)
}

func (h *usersHandlers) VerifyLoginWithRecoveryCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	verifyLoginRequest := &payloads.LoginWithRecoveryCodeRequest{}
	if err = util.ReadRequest(ctx, r, verifyLoginRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	var loginResponse *payloads.LoginResponse
	if loginResponse, err = h.controller.VerifyLoginWithRecoveryCode(ctx, verifyLoginRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, loginResponse)
}

func (h *usersHandlers) Begin2faSetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	currentUser, ok := ctx.Value(middleware.CurrentUserKey).(*models.User)
	if !ok {
		h.logger.HttpError(r, http_errors.NewHttpError(http.StatusInternalServerError, "Error while parsing context", nil))
		base.SendErrorResponse(w, http_errors.NewHttpError(http.StatusInternalServerError, "Invalid credentials", nil))
		return
	}

	var base64QrCode string
	if base64QrCode, err = h.controller.Begin2faSetupSession(ctx, currentUser); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	response := &payloads.Begin2faSetupResponse{
		Base64QrCode: base64QrCode,
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, response)
}

func (h *usersHandlers) Complete2faSetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	currentUser, ok := ctx.Value(middleware.CurrentUserKey).(*models.User)
	if !ok {
		h.logger.HttpError(r, http_errors.NewHttpError(http.StatusInternalServerError, "Error while parsing context", nil))
		base.SendErrorResponse(w, http_errors.NewHttpError(http.StatusInternalServerError, "Invalid credentials.", nil))
		return
	}

	var err error
	complete2faSetupRequest := &payloads.Complete2faSetupRequest{}
	if err = util.ReadRequest(ctx, r, complete2faSetupRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	var recoveryCodes []*models.RecoveryCode
	if recoveryCodes, err = h.controller.Complete2faSetup(ctx, complete2faSetupRequest, currentUser); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	codes := make([]string, 0, len(recoveryCodes))
	for _, code := range recoveryCodes {
		codes = append(codes, code.Code)
	}

	response := &payloads.Complete2faSetupResponse{
		RecoveryCodes: codes,
	}
	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, response)
}

func (h *usersHandlers) Disable2fa(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, ok := ctx.Value(middleware.CurrentUserKey).(*models.User)
	if !ok {
		h.logger.HttpError(r, http_errors.NewHttpError(http.StatusInternalServerError, "Error while parsing context", nil))
		base.SendErrorResponse(w, http_errors.NewHttpError(http.StatusInternalServerError, "Invalid credentials.", nil))
		return
	}

	var err error
	disable2faRequest := &payloads.Disable2faRequest{}
	if err = util.ReadRequest(ctx, r, disable2faRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	if err = h.controller.Disable2fa(ctx, currentUser, disable2faRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, "Successfully disabled two factor auth.")
}

func (h *usersHandlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser, ok := ctx.Value(middleware.CurrentUserKey).(*models.User)
	if !ok {
		h.logger.HttpError(r, http_errors.NewHttpError(http.StatusInternalServerError, "Error while parsing context", nil))
		base.SendErrorResponse(w, http_errors.NewHttpError(http.StatusInternalServerError, "Invalid credentials.", nil))
	}

	var err error
	updatePasswordRequest := &payloads.UpdatePasswordRequest{}
	if err = util.ReadRequest(ctx, r, updatePasswordRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	if err = h.controller.UpdatePassword(ctx, currentUser, updatePasswordRequest); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, "Successfully updated password.")
}

func (h *usersHandlers) RegenerateRecoveryCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	currentUser, ok := ctx.Value(middleware.CurrentUserKey).(*models.User)
	if !ok {
		h.logger.HttpError(r, http_errors.NewHttpError(http.StatusInternalServerError, "Error while parsing context", nil))
		base.SendErrorResponse(w, http_errors.NewHttpError(http.StatusInternalServerError, "Invalid credentials.", nil))
		return
	}

	var recoveryCodes []*models.RecoveryCode
	if recoveryCodes, err = h.controller.RegenerateRecoveryCodes(ctx, currentUser); err != nil {
		h.logger.HttpError(r, err)
		base.SendErrorResponse(w, err)
		return
	}

	codes := make([]string, 0, len(recoveryCodes))
	for _, code := range recoveryCodes {
		codes = append(codes, code.Code)
	}

	response := &payloads.RegenerateRecoveryCodesResponse{
		RecoveryCodes: codes,
	}

	h.logger.HttpSuccess(r)
	base.SendSuccessResponse(ctx, w, response)
}
