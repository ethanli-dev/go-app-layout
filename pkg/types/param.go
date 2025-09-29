/*
Copyright © 2025 lixw
*/
package types

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const MaxPageSize = 100

type PageRequest struct {
	Current  int `json:"current"`  // 当前页码，从1开始
	PageSize int `json:"pageSize"` // 每页大小
}

func (p *PageRequest) Offset() int {
	page := p.PageNo()
	size := p.Limit()
	return (page - 1) * size
}

func (p *PageRequest) Limit() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > MaxPageSize {
		return MaxPageSize
	}
	return p.PageSize
}

func (p *PageRequest) PageNo() int {
	if p.Current <= 0 {
		return 1
	}
	return p.Current
}

type PageData[T any] struct {
	List     []T   `json:"list"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

func NewPageData[T any](list []T, total int64, page, size int) *PageData[T] {
	return &PageData[T]{
		List:     list,
		Page:     page,
		PageSize: size,
		Total:    total,
	}
}

func EmptyPageData[T any](page, size int) *PageData[T] {
	return &PageData[T]{
		List:     make([]T, 0),
		Page:     page,
		PageSize: size,
		Total:    0,
	}
}

type Response[T any] struct {
	Code      int32  `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

func NewResponse[T any](code int32, message string, data T) *Response[T] {
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
	response := NewResponse[any](http.StatusOK, message, nil)
	ctx.JSON(http.StatusOK, response)
}

func SuccessWithData[T any](ctx *gin.Context, data T, messages ...string) {
	message := "success"
	if len(messages) > 0 {
		message = messages[0]
	}
	response := NewResponse(http.StatusOK, message, data)
	ctx.JSON(http.StatusOK, response)
}
