package http_errors

import (
	"fmt"
	"net/http"
)

type Error interface {
	Error() string
	Reasons() string
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

func (e *HttpError) Reasons() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.status, e.message, e.reasons)
}

func (e *HttpError) Message() string {
	return e.message
}

func (e *HttpError) Status() int {
	return e.status
}

func NewInternalServerError(reasons interface{}) Error {
	result := &HttpError{
		status:  http.StatusInternalServerError,
		message: "Internal server error.",
		reasons: reasons,
	}
	return result
}

func NewBadGatewayError(reasons interface{}) Error {
	result := &HttpError{
		status:  http.StatusBadGateway,
		message: "Bad gateway.",
		reasons: reasons,
	}
	return result
}

func NewBadRequestError(reasons interface{}) Error {
	result := &HttpError{
		status:  http.StatusBadRequest,
		message: "Bad request.",
		reasons: reasons,
	}
	return result
}

func NewForbiddenError(reasons interface{}) Error {
	result := &HttpError{
		status:  http.StatusForbidden,
		message: "Forbidden error.",
		reasons: reasons,
	}
	return result
}

func NewNotFoundError(reasons interface{}) Error {
	result := &HttpError{
		status:  http.StatusNotFound,
		message: "Not found.",
		reasons: reasons,
	}
	return result
}

func InvalidCredentialsError(reasons interface{}) Error {
	result := &HttpError{
		status:  http.StatusBadRequest,
		message: "Invalid credentials.",
		reasons: reasons,
	}
	return result
}
