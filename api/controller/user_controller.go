package controller

import "example/dashboard/api/store"

type UserController struct {
	UserStore store.UserStore
}

func (s UserController) CreateUser(email string, password string, firstName string, lastName string) error {
	return s.UserStore.CreateUser(email, password, firstName, lastName)
}
