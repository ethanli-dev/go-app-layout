/*
Copyright © 2025 lixw
*/
package handler

import (
	"github.com/ethanli-dev/go-app-layout/api"
	v1 "github.com/ethanli-dev/go-app-layout/api/v1"
	"github.com/ethanli-dev/go-app-layout/internal/service"
	"github.com/ethanli-dev/go-app-layout/pkg/errorx"
	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	tenantSrv *service.TenantService
}

func NewTenantHandler(tenantSrv *service.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantSrv: tenantSrv,
	}
}

// Create tenant
// @Summary Create tenant
// @Tags tenant
// @Accept json
// @Produce json
// @Param req body v1.TenantRequest true "Create tenant request"
// @Success 200 {object} api.Response[any] "Create tenant response"
// @Router /tenant/create [post]
func (th *TenantHandler) Create(ctx *gin.Context) {
	var req v1.TenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Failure(ctx, errorx.New(errorx.ErrCodeBadRequest, "failed to parse request parameters"))
		return
	}
	create, err := th.tenantSrv.Create(ctx, &req)
	if err != nil {
		api.Failure(ctx, err)
		return
	}
	api.SuccessWithData(ctx, create)
}
