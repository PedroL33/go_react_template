package middleware

import (
	"context"
	"example/dashboard/api/httpcomm"
	"example/dashboard/api/models"
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
			m.logger.HttpError(r, httpcomm.NewInternalServerError(err, "Authorization header not found."))
			httpcomm.SendErrorResponse(w, httpcomm.NewInternalServerError(err, "Invalid credentials."))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			m.logger.HttpError(r, httpcomm.NewInternalServerError(err, "Missing token."))
			httpcomm.SendErrorResponse(w, httpcomm.NewInternalServerError(err, "Invalid credentials."))
			return
		}

		var payload map[string]interface{}
		if payload, err = util.VerifyToken(m.conf, token); err != nil {
			m.logger.HttpError(r, httpcomm.NewInternalServerError(err, "Error while verifying token."))
			httpcomm.SendErrorResponse(w, httpcomm.NewInternalServerError(err, "Invalid credentials."))
			return
		}

		var currentUser *models.User
		if username, ok := payload["Username"].(string); ok {
			if currentUser, err = m.userStore.GetUserByUsername(r.Context(), username); err != nil {
				m.logger.HttpError(r, httpcomm.NewInternalServerError(err, "Invalid token payload."))
				httpcomm.SendErrorResponse(w, httpcomm.NewInternalServerError(err, "Invalid credentials."))
				return
			}
		}
		ctx := context.WithValue(r.Context(), CurrentUserKey, currentUser)
		next(w, r.WithContext(ctx))

	})
}
