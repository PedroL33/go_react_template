package util

import (
	"context"
	http_errors "example/dashboard/errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate
var trans ut.Translator

func init() {
	validate = validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	err := enTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("Failed to register default translations:", err)
	}
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	fmt.Println("Validator initialized")
}

// Validate struct fields
func ValidateStruct(ctx context.Context, s interface{}) error {

	err := validate.StructCtx(ctx, s)

	if err != nil {
		translatedErrors := make(map[string]string)
		// Loop through validation errors and collect the translated errors
		for _, e := range err.(validator.ValidationErrors) {
			// Translate the error and append to the slice
			translatedErrors[e.Field()] = e.Translate(trans)
		}

		if len(translatedErrors) > 0 {
			return http_errors.NewHttpError(422, "Validation failed.", translatedErrors)
		}
	}

	return nil
}
