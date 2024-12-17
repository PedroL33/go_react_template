package util

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

func ReadRequest(ctx context.Context, r *http.Request, request interface{}) error {

	err := DecodeInto(r, request)
	if err != nil {
		return errors.Wrap(err, "util.ReadRequest")
	}

	return validate.StructCtx(ctx, request)
}
