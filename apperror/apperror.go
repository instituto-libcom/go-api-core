package apperror

import "net/http"

// AppError represents a standard error for the API
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Cause      error
}

// Error implements the standard error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap allows for error unwrapping
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Wrap preserves the original error
func (e *AppError) Wrap(cause error) *AppError {
	e.Cause = cause
	return e
}

// Semantic Constructors

// NewBadRequest returns a 400 error
func NewBadRequest(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusBadRequest}
}

// NewNotFound returns a 404 error
func NewNotFound(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusNotFound}
}

// NewConflict returns a 409 error
func NewConflict(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusConflict}
}

// NewForbidden returns a 403 error
func NewForbidden(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusForbidden}
}

// NewUnprocessable returns a 422 error
func NewUnprocessable(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusUnprocessableEntity}
}

// NewUnavailable returns a 503 error
func NewUnavailable(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusServiceUnavailable}
}

// NewInternal returns a 500 error
func NewInternal(code, message string) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: http.StatusInternalServerError}
}
