package base

import (
	"context"
	"errors"
	http_errors "example/dashboard/errors"
	"example/dashboard/util"
	"net/http"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

func SendErrorResponse(w http.ResponseWriter, err error) {

	// var ve validator.ValidationErrors

	// if errors.As(err, &ve) {
	// 	out := make(map[string]string)
	// 	for _, fe := range ve {
	// 		out[fe.Field()] = fe.Error()
	// 	}
	// 	errorResponse := &ErrorResponse{
	// 		Status:  http.StatusUnprocessableEntity,
	// 		Message: "Failed validation",
	// 		Error:   out,
	// 	}
	// 	util.Encode(w, http.StatusUnprocessableEntity, errorResponse)
	// }

	var httpError *http_errors.HttpError

	if errors.As(err, &httpError) {
		errorResponse := &ErrorResponse{
			Status:  httpError.Status(),
			Message: httpError.Message(),
			Error:   httpError.Reasons(),
		}

		util.Encode(w, httpError.Status(), errorResponse)
	}
}
