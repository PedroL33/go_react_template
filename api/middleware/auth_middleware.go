package middleware

import (
	"context"
	"example/dashboard/api/base"
	"example/dashboard/api/models"
	http_errors "example/dashboard/errors"
	"example/dashboard/util"
	"net/http"
	"strings"
)

type contextKey string

const CurrentUserKey contextKey = "currentUser"

func (m *middleWareManager) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.HttpError(r, http_errors.NewInternalServerError(err, "Authorization header not found."))
			base.SendErrorResponse(w, http_errors.NewInternalServerError(err, "Invalid credentials."))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer: ")
		if token == "" {
			m.logger.HttpError(r, http_errors.NewInternalServerError(err, "Missing token."))
			base.SendErrorResponse(w, http_errors.NewInternalServerError(err, "Invalid credentials."))
			return
		}

		var payload map[string]interface{}
		if payload, err = util.VerifyToken(m.conf, token); err != nil {
			m.logger.HttpError(r, http_errors.NewInternalServerError(err, "Error while verifying token."))
			base.SendErrorResponse(w, http_errors.NewInternalServerError(err, "Invalid credentials."))
			return
		}

		var currentUser *models.User
		if userEmail, ok := payload["Email"].(string); ok {
			if currentUser, err = m.userStore.GetUserByEmail(r.Context(), userEmail, nil); err != nil {
				m.logger.HttpError(r, http_errors.NewInternalServerError(err, "Invalid token payload."))
				base.SendErrorResponse(w, http_errors.NewInternalServerError(err, "Invalid credentials."))
				return
			}
		}
		ctx := context.WithValue(r.Context(), CurrentUserKey, currentUser)
		next(w, r.WithContext(ctx))

	})
}
