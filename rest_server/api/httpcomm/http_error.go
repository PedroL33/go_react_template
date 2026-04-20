package httpcomm

import (
	"net/http"
)

type HttpError struct {
	Status  int
	Message string // user-facing
	Err     error  // underlying cause, for logging only
	Reasons any    // optional structured payload (e.g. validation map)
}

func (e *HttpError) Error() string { return e.Message }

// Unwrap lets errors.Is / errors.As traverse into the cause.
func (e *HttpError) Unwrap() error { return e.Err }

func (e *HttpError) Trace() string { return e.Err.Error() }

// Constructors — take the cause as an error, not interface{}.
func NewInternalServerError(cause error, msg string) *HttpError {
	return &HttpError{Status: http.StatusInternalServerError, Message: msg, Err: cause}
}
func NewBadRequestError(cause error, msg string) *HttpError {
	return &HttpError{Status: http.StatusBadRequest, Message: msg, Err: cause}
}
func NewUnauthorizedError(cause error, msg string) *HttpError {
	return &HttpError{Status: http.StatusUnauthorized, Message: msg, Err: cause}
}
func NewForbiddenError(cause error, msg string) *HttpError {
	return &HttpError{Status: http.StatusForbidden, Message: msg, Err: cause}
}
func NewUnprocessableEntityError(reasons any, msg string) *HttpError {
	return &HttpError{Status: http.StatusUnprocessableEntity, Message: msg, Reasons: reasons}
}
