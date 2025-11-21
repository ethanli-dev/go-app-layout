/*
Copyright © 2025 lixw
*/
package server

import (
	"context"
	"log/slog"

	"github.com/ethanli-dev/go-app-layout/internal/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	tenantHandler *handler.TenantHandler
}

func New(tenantHandler *handler.TenantHandler) *Server {
	return &Server{
		tenantHandler: tenantHandler,
	}
}

func (s *Server) Routes(group *gin.RouterGroup) {
	authGroup := group.Group("/tenant")
	{
		authGroup.POST("/create", s.tenantHandler.Create)
	}
}

func (s *Server) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "starting internal server")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "stopping internal server")
	return nil
}
