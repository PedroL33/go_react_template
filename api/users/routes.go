package users

import (
	"github.com/gorilla/mux"
)

func MapRoutes(h Handler, router *mux.Router) {
	router.HandleFunc("/user", h.Create).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
}
