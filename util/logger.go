package util

import (
	"bytes"
	"errors"
	http_errors "example/dashboard/errors"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Logger interface {
	Error(err error, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	HttpError(r *http.Request, err error)
	HttpSuccess(r *http.Request)
}

type logger struct {
	lgr *slog.Logger
}

func NewLogger() Logger {
	return &logger{
		lgr: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

func (l *logger) Error(err error, args ...any) {
	l.lgr.Error(err.Error(), args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.lgr.Info("INFO: "+msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.lgr.Warn("WARN: "+msg, args...)
}

func (l *logger) HttpError(r *http.Request, err error) {
	var httpError *http_errors.HttpError
	var errorMessage string
	if errors.As(err, &httpError) {
		errorMessage = httpError.Reasons()
	} else {
		l.Warn("Recieved an error that is not an HttpError.")
		errorMessage = err.Error()
	}
	l.lgr.Error(
		"ERROR: "+errorMessage,
		slog.String("URI", r.RequestURI),
		slog.String("Method", r.Method),
		slog.String("Address", r.RemoteAddr),
		slog.String("Body", GetRequestBody(r.Body)),
	)
}

func (l *logger) HttpSuccess(r *http.Request) {
	l.Info(
		"Successful request",
		slog.String("URI", r.RequestURI),
		slog.String("Method", r.Method),
		slog.String("Address", r.RemoteAddr),
		slog.String("Body", GetRequestBody(r.Body)),
	)
}

func GetRequestBody(body io.ReadCloser) string {
	// Read the body bytes
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return "Failed to read request body"
	}

	// Restore the body so it can be read again
	defer body.Close() // Close the original body
	io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Convert the body to a string
	bodyString := string(bodyBytes)

	return bodyString

}
