package errors

type AppError struct{
	Code string `json:"code"`
	Message string `json:"message"`
}

func(c * AppError) Error() string{
	return c.Message
}

func NewNotFound(message string) error{
	return &AppError{Code:"Not_Found",Message: message}
}

func NewConflict(message string) error {
	return &AppError{Code: "CONFLICT",Message: message}
}

func NewUnauthorized(message string) error{
	return &AppError{Code:"UNAUTHORIZED", Message:message}
}

func NewValidation(message string) error {
	return &AppError{Code: "VALIDATION_ERROR",Message: message}
}

func NewInternal(message string) error {
	return &AppError{Code: "INTERNAL_ERROR",Message: message}
}