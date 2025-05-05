package validate

import (
	"fmt"

	"github.com/go-playground/validator"
)

func ValidateRequestBody(request any) map[string]any {
	var errs map[string]any

	err := validator.New().Struct(request)
	if err != nil {
		var errors []string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				errors = append(errors, getValidationMsg(fieldErr))
			}
		}

		errs = map[string]any{"errors": errors}
	}

	return errs
}

func getValidationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	}

	return "validation error"
}
