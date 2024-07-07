package utils

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func generateValidationParamErrorMsg(paramName string) string {
	return fmt.Sprintf("invalid param %s", paramName)
}

func ValidateURLParamUUID(r *http.Request, paramName string, defaultValue ...uuid.UUID) uuid.UUID {
	param := chi.URLParam(r, paramName)

	uuid, err := uuid.Parse(param)
	if err != nil {
		if len(defaultValue) > 0 {
			uuid = defaultValue[0]
		} else {
			PanicIfError(CustomErrorWithTrace(err, generateValidationParamErrorMsg(paramName), 400))
		}
	}

	return uuid
}
