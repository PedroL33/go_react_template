package util

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errors.Wrap(err, "util.Encode")
	}
	return nil
}
