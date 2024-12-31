package util

import (
	"context"
	http_errors "example/dashboard/errors"
	"net/http"

	"github.com/pkg/errors"
)

func ReadRequest(ctx context.Context, r *http.Request, request interface{}) error {

	err := DecodeInto(r, request)
	if err != nil {
		// This is only called in the handler level so we want ot return an http error
		return http_errors.NewInternalServerError(errors.Wrap(err, "util.ReadRequest"), "Error while reading request.")
	}

	return ValidateStruct(ctx, request)
}
