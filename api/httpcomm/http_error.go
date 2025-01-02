package httpcomm

import (
	"fmt"
	"net/http"
)

type Error interface {
	Error() string
	Reasons() interface{}
	LoggerMessage() string
	Message() string
	Status() int
}

type HttpError struct {
	status  int
	message string
	reasons interface{}
}

func NewHttpError(code int, msg string, reasons interface{}) *HttpError {
	return &HttpError{
		status:  code,
		message: msg,
		reasons: reasons,
	}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s", e.status, e.message)
}

func (e *HttpError) LoggerMessage() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.status, e.message, e.reasons)
}

func (e *HttpError) Reasons() interface{} {
	return e.reasons
}

func (e *HttpError) Message() string {
	return e.message
}

func (e *HttpError) Status() int {
	return e.status
}

func NewInternalServerError(reasons interface{}, msg string) Error {
	result := &HttpError{
		status:  http.StatusInternalServerError,
		message: msg,
		reasons: reasons,
	}
	return result
}

func NewBadGatewayError(reasons interface{}, msg string) Error {
	result := &HttpError{
		status:  http.StatusBadGateway,
		message: msg,
		reasons: reasons,
	}
	return result
}

func NewBadRequestError(reasons interface{}, msg string) Error {
	result := &HttpError{
		status:  http.StatusBadRequest,
		message: msg,
		reasons: reasons,
	}
	return result
}

func NewForbiddenError(reasons interface{}, msg string) Error {
	result := &HttpError{
		status:  http.StatusForbidden,
		message: msg,
		reasons: reasons,
	}
	return result
}

func NewNotFoundError(reasons interface{}, msg string) Error {
	result := &HttpError{
		status:  http.StatusNotFound,
		message: msg,
		reasons: reasons,
	}
	return result
}

func NewExpiredSessionError(reasons interface{}, msg string) Error {
	result := &HttpError{
		status:  http.StatusGone,
		message: msg,
		reasons: reasons,
	}
	return result
}
