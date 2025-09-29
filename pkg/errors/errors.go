/*
Copyright © 2025 lixw
*/
package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode int

// System error codes
const (
	// Common error codes (10000-10999)
	ErrBadRequest         ErrorCode = 10000
	ErrUnauthorized       ErrorCode = 10001
	ErrForbidden          ErrorCode = 10002
	ErrNotFound           ErrorCode = 10003
	ErrMethodNotAllowed   ErrorCode = 10004
	ErrConflict           ErrorCode = 10005
	ErrTooManyRequests    ErrorCode = 10006
	ErrInternalServer     ErrorCode = 10007
	ErrServiceUnavailable ErrorCode = 10008
	ErrTimeout            ErrorCode = 10009
	ErrValidation         ErrorCode = 10010
)

type AppError struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	HTTPCode int       `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("error code: %d, error message: %s", e.Code, e.Message)
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:     ErrBadRequest,
		Message:  message,
		HTTPCode: http.StatusBadRequest,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:     ErrUnauthorized,
		Message:  message,
		HTTPCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:     ErrForbidden,
		Message:  message,
		HTTPCode: http.StatusForbidden,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:     ErrNotFound,
		Message:  message,
		HTTPCode: http.StatusNotFound,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:     ErrConflict,
		Message:  message,
		HTTPCode: http.StatusConflict,
	}
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(message string) *AppError {
	if message == "" {
		message = "服务器内部错误"
	}
	return &AppError{
		Code:     ErrInternalServer,
		Message:  message,
		HTTPCode: http.StatusInternalServerError,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:     ErrValidation,
		Message:  message,
		HTTPCode: http.StatusBadRequest,
	}
}

// IsAppError checks if the error is an AppError type
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
