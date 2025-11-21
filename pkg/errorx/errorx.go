/*
Copyright © 2025 lixw
*/
package errorx

import (
	"errors"
	"fmt"
)

const (
	// ErrCodeSuccess 标识请求成功
	ErrCodeSuccess = 200
	// ErrCodeBadRequest 表示请求参数错误（对应HTTP 400）
	ErrCodeBadRequest = 10000
	// ErrCodeUnauthorized 表示未授权（对应HTTP 401）
	ErrCodeUnauthorized = 10001
	// ErrCodeForbidden 表示权限不足（对应HTTP 403）
	ErrCodeForbidden = 10002
	// ErrCodeNotFound 表示资源不存在（对应HTTP 404）
	ErrCodeNotFound = 10003
	// ErrCodeMethodNotAllowed 表示方法不支持（对应HTTP 405）
	ErrCodeMethodNotAllowed = 10004
	// ErrCodeConflict 表示资源冲突（对应HTTP 409）
	ErrCodeConflict = 10005
	// ErrCodeTooManyRequests 表示请求过于频繁（对应HTTP 429）
	ErrCodeTooManyRequests = 10006
	// ErrCodeInternalServer 表示服务器内部错误（对应HTTP 500）
	ErrCodeInternalServer = 10007
	// ErrCodeServiceUnavailable 表示服务不可用（对应HTTP 503）
	ErrCodeServiceUnavailable = 10008
	// ErrCodeTimeout 表示请求超时（对应HTTP 504）
	ErrCodeTimeout = 10009
	// ErrCodeValidation 表示数据校验失败
	ErrCodeValidation = 10010
	// ErrCodePermissionDenied 表示权限拒绝（更细化的权限错误）
	ErrCodePermissionDenied = 10011
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

// IsCode 判断错误链中是否包含指定错误码
func IsCode(err error, code int) bool {
	for err != nil {
		var e *WrappedError
		if errors.As(err, &e) {
			if e.Code == code {
				return true
			}
		}
		err = errors.Unwrap(err)
	}
	return false
}

// CodeOf 提取错误链中第一个 WrappedError 的 code
func CodeOf(err error) int {
	for err != nil {
		var e *WrappedError
		if errors.As(err, &e) {
			return e.Code
		}
		err = errors.Unwrap(err)
	}
	return ErrCodeInternalServer
}

// MessageOf 提取错误链中第一个 WrappedError 的 message
func MessageOf(err error) string {
	for err != nil {
		var e *WrappedError
		if errors.As(err, &e) {
			return e.Message
		}
		err = errors.Unwrap(err)
	}
	return ""
}
