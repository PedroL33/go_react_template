package httpcomm

import (
	"context"
	"example/dashboard/util"
	"example/dashboard/validator"
	"net/http"

	"github.com/pkg/errors"
)

func ReadRequest(ctx context.Context, r *http.Request, request interface{}) error {

	err := util.DecodeInto(r, request)
	if err != nil {
		// This is only called in the handler level so we want ot return an http error
		return NewInternalServerError(errors.Wrap(err, "util.ReadRequest"), "Error while reading request.")
	}

	errors := validator.ValidateStruct(ctx, request)

	if errors != nil {
		return NewHttpError(422, "Validation failed.", errors)
	}

	return nil
}
