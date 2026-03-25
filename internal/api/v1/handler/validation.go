package handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(input any) map[string]string {
	err := validate.Struct(input)
	if err == nil {
		return nil
	}
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "Unvalid"}
	}

	result := make(map[string]string, len(validationErrs))
	for _, fieldErr := range validationErrs {
		field := strings.ToLower(fieldErr.Field())
		result[field] = mapValidationMessage(fieldErr)
	}
	return result
}

func mapValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", strings.ToLower(e.Field()))
	case "email":
		return fmt.Sprintf("%s must be a valid email", strings.ToLower(e.Field()))
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", strings.ToLower(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", strings.ToLower(e.Field()), e.Param())
	default:
		return fmt.Sprintf("%s không hợp lệ", strings.ToLower(e.Field()))
	}
}
