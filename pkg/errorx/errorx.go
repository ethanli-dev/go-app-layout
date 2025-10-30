/*
Copyright © 2025 lixw
*/
package errorx

import (
	"errors"
	"fmt"
)

const (
	// Common error codes (10000-10999)
	ErrBadRequest         = 10000
	ErrUnauthorized       = 10001
	ErrForbidden          = 10002
	ErrNotFound           = 10003
	ErrMethodNotAllowed   = 10004
	ErrConflict           = 10005
	ErrTooManyRequests    = 10006
	ErrInternalServer     = 10007
	ErrServiceUnavailable = 10008
	ErrTimeout            = 10009
	ErrValidation         = 10010
	ErrPermissionDenied   = 10011
)

type WrappedError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

// New 创建一个新的基础错误
func New(code int, format string, args ...any) *WrappedError {
	return &WrappedError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap 包装已有错误
func Wrap(err error, code int, format string, args ...any) *WrappedError {
	return &WrappedError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// Error 实现error接口
func (e *WrappedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 兼容 errors.Unwrap
func (e *WrappedError) Unwrap() error {
	return e.Cause
}

// JSON 返回结构化错误信息
func (e *WrappedError) JSON() map[string]any {
	return map[string]any{
		"code":    e.Code,
		"message": e.Message,
	}
}

// Is 判断错误链中是否包含指定错误码
func Is(err error, code int) bool {
	if err == nil {
		return false
	}
	var e *WrappedError
	if errors.As(err, &e) {
		if e.Code == code {
			return true
		}
	}
	return Is(errors.Unwrap(err), code)
}

// CodeOf 提取错误链中第一个 WrappedError 的 code
func CodeOf(err error) int {
	if err == nil {
		return -1
	}
	var e *WrappedError
	if errors.As(err, &e) {
		return e.Code
	}
	return CodeOf(errors.Unwrap(err))
}

// MessageOf 提取错误链中第一个 WrappedError 的 message
func MessageOf(err error) string {
	if err == nil {
		return ""
	}
	var e *WrappedError
	if errors.As(err, &e) {
		return e.Message
	}
	return MessageOf(errors.Unwrap(err))
}
