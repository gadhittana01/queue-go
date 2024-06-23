package utils

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Message    string
	StatusCode int
	ErrorCode  int
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("app error: status code %d, message %s", ae.StatusCode, ae.Message)
}

type ValidationError struct {
	Message string
	Field   string
	Tag     string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error: message %s", ve.Message)
}

type ValidationErrors struct {
	Errors     []ValidationError
	StatusCode int
}

func CustomError(message string, statusCode int) error {
	return fmt.Errorf("%s|%s<->%d|", message, message, statusCode)
}

func CustomErrorWithTrace(err error, message string, statusCode int) error {
	return fmt.Errorf("%s|%s<->%d|", err.Error(), message, statusCode)
}

func PanicIfError(err error) {
	if err != nil {
		customError := strings.Split(err.Error(), "<->")
		message := customError[0]
		statusCode := 500

		if len(customError) > 1 {
			statusCode, _ = strconv.Atoi(strings.Split(customError[1], "|")[0])
		}

		appErr := AppError{
			Message:    message,
			StatusCode: statusCode,
		}
		panic(appErr)
	}
}

func PanicIfAppError(err error, message string, statusCode int) {
	if err != nil {
		customErr := CustomErrorWithTrace(err, message, statusCode)
		PanicIfError(customErr)
	}
}

func PanicAppError(message string, statusCode int) {
	customErr := CustomError(message, statusCode)
	PanicIfError(customErr)
}

func PanicValidationError(errors []ValidationError, statusCode int) {
	validationErrors := ValidationErrors{
		Errors:     errors,
		StatusCode: statusCode,
	}
	panic(validationErrors)
}

func ValidateStruct(data interface{}) {
	var validationErrors []ValidationError
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	errorValidate := validate.Struct(data)

	if errorValidate != nil {
		for _, err := range errorValidate.(validator.ValidationErrors) {
			var validationError ValidationError
			validationError.Message = strings.Split(err.Error(), "Error:")[1]
			validationError.Field = err.Field()
			validationError.Tag = err.Tag()
			validationErrors = append(validationErrors, validationError)
		}
		PanicValidationError(validationErrors, 400)
	}
}

func ValidateBodyPayload[T any](body io.ReadCloser, output *T) T {
	err := JSONiter().NewDecoder(body).Decode(output)
	PanicIfAppError(err, "failed when decode body payload", 400)

	ValidateStruct(output)
	return *output
}
