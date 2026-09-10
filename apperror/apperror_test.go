package apperror

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	err := &AppError{Message: "Some error occurred"}
	if err.Error() != "Some error occurred" {
		t.Errorf("Expected 'Some error occurred', got '%s'", err.Error())
	}

	cause := errors.New("root cause")
	errWithCause := err.Wrap(cause)

	expected := "Some error occurred: root cause"
	if errWithCause.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, errWithCause.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	err := &AppError{Message: "wrapped error"}
	err.Wrap(cause)

	if !errors.Is(err, cause) {
		t.Errorf("errors.Is should return true for the cause")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("Unwrap did not return the expected cause")
	}
}

func TestSemanticConstructors(t *testing.T) {
	tests := []struct {
		name       string
		err        *AppError
		statusCode int
	}{
		{"BadRequest", NewBadRequest("BAD_REQ", "Bad Request"), http.StatusBadRequest},
		{"NotFound", NewNotFound("NOT_FOUND", "Not Found"), http.StatusNotFound},
		{"Conflict", NewConflict("CONFLICT", "Conflict"), http.StatusConflict},
		{"Forbidden", NewForbidden("FORBIDDEN", "Forbidden"), http.StatusForbidden},
		{"Unprocessable", NewUnprocessable("UNPROCESSABLE", "Unprocessable"), http.StatusUnprocessableEntity},
		{"Unavailable", NewUnavailable("UNAVAILABLE", "Unavailable"), http.StatusServiceUnavailable},
		{"Internal", NewInternal("INTERNAL", "Internal"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.StatusCode != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, tt.err.StatusCode)
			}
		})
	}
}
