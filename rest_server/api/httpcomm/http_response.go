package httpcomm

import (
	"context"
	"errors"
	"example/template/rest_server/util"
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
	var httpError *HttpError
	if !errors.As(err, &httpError) {
		httpError = NewInternalServerError(err, "Internal server error.")
	}
	util.Encode(w, httpError.Status, ErrorResponse{
		Status:  httpError.Status,
		Message: httpError.Message,
		Error:   httpError.Reasons,
	})
}
