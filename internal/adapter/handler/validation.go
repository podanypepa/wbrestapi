package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorResponse represents a single field validation error
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors converts validator.ValidationErrors into a slice of ValidationErrorResponse
func FormatValidationErrors(err error) []ValidationErrorResponse {
	var errors []ValidationErrorResponse
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			errors = append(errors, ValidationErrorResponse{
				Field:   fe.Field(),
				Message: formatFieldError(fe),
			})
		}
	}
	return errors
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "uuid":
		return "must be a valid UUID"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", fe.Param())
	default:
		return fmt.Sprintf("failed on tag %s", fe.Tag())
	}
}
