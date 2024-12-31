package middleware

import (
	"example/dashboard/api/users"
	"example/dashboard/config"
	"example/dashboard/util"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MiddleWareManager interface {
	CompileMiddlewares(h http.HandlerFunc, args ...Middleware) http.HandlerFunc
	Cors(next http.Handler) http.Handler
	Auth(next http.HandlerFunc) http.HandlerFunc
}

type middleWareManager struct {
	conf      *config.Config
	logger    util.Logger
	userStore users.Store
}

func NewMiddlewareManager(conf *config.Config, logger util.Logger, userStore users.Store) MiddleWareManager {
	return &middleWareManager{
		conf:      conf,
		logger:    logger,
		userStore: userStore,
	}
}

func (m *middleWareManager) CompileMiddlewares(h http.HandlerFunc, args ...Middleware) http.HandlerFunc {
	if len(args) < 1 {
		return h
	}

	wrapped := h

	for i := len(args) - 1; i >= 0; i-- {
		wrapped = args[i](wrapped)
	}

	return wrapped
}
