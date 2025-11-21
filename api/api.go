/*
Copyright © 2025 lixw
*/
package api

import (
	"net/http"
	"time"

	"github.com/ethanli-dev/go-app-layout/pkg/errorx"
	"github.com/gin-gonic/gin"
)

type PageRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

func (p *PageRequest) Offset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

func (p *PageRequest) Limit() int {
	return p.GetPageSize()
}

type PageResult[T any] struct {
	List     []T `json:"list"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func NewPageResult[T any](list []T, page, size, total int) *PageResult[T] {
	return &PageResult[T]{
		List:     list,
		Page:     page,
		PageSize: size,
		Total:    total,
	}
}
func EmptyPageResult[T any](page, size int) *PageResult[T] {
	return &PageResult[T]{
		List:     []T{},
		Page:     page,
		PageSize: size,
		Total:    0,
	}
}

type Response[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

func NewResponse[T any](code int, message string, data T) *Response[T] {
	return &Response[T]{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

func Success(ctx *gin.Context, messages ...string) {
	message := "success"
	if len(messages) > 0 {
		message = messages[0]
	}
	ctx.JSON(http.StatusOK, NewResponse[any](errorx.ErrCodeSuccess, message, nil))
}

func SuccessWithData[T any](ctx *gin.Context, data T, messages ...string) {
	message := "success"
	if len(messages) > 0 {
		message = messages[0]
	}
	ctx.JSON(http.StatusOK, NewResponse(errorx.ErrCodeSuccess, message, data))
}

func Failure(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, NewResponse[any](errorx.CodeOf(err), errorx.MessageOf(err), nil))
}
