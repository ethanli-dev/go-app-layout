/*
Copyright © 2025 lixw
*/
package repository

import (
	"context"
	"errors"

	"github.com/ethanli-dev/go-app-layout/internal/model"
	"gorm.io/gorm"
)

var ErrTenantNotFound = errors.New("tenant not found")

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

func (tr *TenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	return tr.db.WithContext(ctx).Create(tenant).Error
}

func (tr *TenantRepository) GetById(ctx context.Context, id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := tr.db.WithContext(ctx).First(&tenant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (tr *TenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	return tr.db.WithContext(ctx).Model(&model.Tenant{}).Where("id = ?", tenant.ID).Updates(tenant).Error
}

func (tr *TenantRepository) Delete(ctx context.Context, id uint) error {
	return tr.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Tenant{}).Error
}
