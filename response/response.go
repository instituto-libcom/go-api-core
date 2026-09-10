package response

import (
	"encoding/json"
	"net/http"

	"github.com/instituto-libcom/go-api-core/pagination"
)

// ErrorDetails defines the standard error structure
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BaseResponse defines the root JSON structure for responses
type BaseResponse[T any] struct {
	Success bool          `json:"success"`
	Data    T             `json:"data"`
	Error   *ErrorDetails `json:"error"`
}

// PaginatedData defines the structure of data in a paginated response
type PaginatedData[T any] struct {
	Items T               `json:"items"`
	Meta  pagination.Meta `json:"meta"`
}

// JSON sends a generic JSON response
func JSON[T any](w http.ResponseWriter, statusCode int, payload BaseResponse[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success sends a single entity success response
func Success[T any](w http.ResponseWriter, statusCode int, data T) {
	JSON(w, statusCode, BaseResponse[T]{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

// Paginated sends a paginated list success response
func Paginated[T any](w http.ResponseWriter, statusCode int, items T, meta pagination.Meta) {
	JSON(w, statusCode, BaseResponse[PaginatedData[T]]{
		Success: true,
		Data: PaginatedData[T]{
			Items: items,
			Meta:  meta,
		},
		Error: nil,
	})
}

// Error sends a standardized error response
func Error(w http.ResponseWriter, statusCode int, code string, message string) {
	JSON(w, statusCode, BaseResponse[any]{
		Success: false,
		Data:    nil,
		Error: &ErrorDetails{
			Code:    code,
			Message: message,
		},
	})
}

