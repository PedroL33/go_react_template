package base

import (
	"context"
	"errors"
	http_errors "example/dashboard/errors"
	"example/dashboard/util"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type SuccessResponse struct {
	Status  int
	Message string
	Data    interface{}
}

func SendSuccessResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	successResponse := &SuccessResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    data,
	}

	util.Encode(w, http.StatusOK, successResponse)
}

type ErrorResponse struct {
	Status  int
	Message string
	Error   interface{}
}

func SendErrorResponse(w http.ResponseWriter, err error) {

	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make(map[string]string)
		for _, fe := range ve {
			out[fe.Field()] = fe.Error()
		}
		errorResponse := &ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed validation",
			Error:   out,
		}
		util.Encode(w, http.StatusBadRequest, errorResponse)
	}

	var httpError *http_errors.HttpError

	if errors.As(err, &httpError) {
		errorResponse := &ErrorResponse{
			Status:  httpError.Status(),
			Message: httpError.Message(),
			Error:   httpError.Error(),
		}

		util.Encode(w, httpError.Status(), errorResponse)
	}
}
