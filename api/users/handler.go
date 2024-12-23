package users

import (
	"net/http"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Begin2faSetup(w http.ResponseWriter, r *http.Request)
	Complete2faSetup(w http.ResponseWriter, r *http.Request)
	Disable2fa(w http.ResponseWriter, r *http.Request)
}
