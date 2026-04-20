package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func DecodeNew[T any](r *http.Request) (T, error) {
	var v T

	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return v, Wrap(err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, Wrap(err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return v, nil
}

func DecodeInto[T any](r *http.Request, value T) error {

	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return Wrap(err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return Wrap(err)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}
