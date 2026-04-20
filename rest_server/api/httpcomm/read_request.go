package httpcomm

import (
	"context"
	"example/template/rest_server/util"
	"example/template/rest_server/validator"
	"net/http"
)

func ReadRequest(ctx context.Context, r *http.Request, request interface{}) error {

	err := util.DecodeInto(r, request)
	if err != nil {
		// This is only called in the handler level so we want ot return an http error
		return NewInternalServerError(util.Wrap(err), "Error while reading request.")
	}

	errors := validator.ValidateStruct(ctx, request)

	if errors != nil {
		return NewUnprocessableEntityError(errors, "Validation failed.")
	}

	return nil
}
