/*
Copyright © 2025 lixw
*/
package router

import (
	"time"

	"github.com/ethanli-dev/go-app-layout/pkg/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) Register(group *gin.RouterGroup) {
	group.GET("/sql", func(ctx *gin.Context) {
		r.db.Exec("select 1")
		types.SuccessWithData(ctx, struct {
			Data int64
		}{
			Data: time.Now().Unix(),
		})
	})
}
