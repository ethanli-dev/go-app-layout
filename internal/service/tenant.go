/*
Copyright © 2025 lixw
*/
package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"
	"log/slog"

	v1 "github.com/ethanli-dev/go-app-layout/api/v1"
	"github.com/ethanli-dev/go-app-layout/internal/model"
	"github.com/ethanli-dev/go-app-layout/internal/repository"
	"github.com/ethanli-dev/go-app-layout/pkg/config"
	"github.com/ethanli-dev/go-app-layout/pkg/errorx"
)

var apiKeySecret = func() []byte {
	return []byte(config.GetString("tenant.aes_key"))
}

type TenantService struct {
	tenantRepo *repository.TenantRepository
}

func NewTenantService(tenantRepo *repository.TenantRepository) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
	}
}

func (tr *TenantService) Create(ctx context.Context, req *v1.TenantRequest) (*model.Tenant, error) {
	if req.Name == "" {
		return nil, errorx.New(errorx.ErrCodeValidation, "name is required")
	}
	tenant := &model.Tenant{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := tr.tenantRepo.Create(ctx, tenant); err != nil {
		slog.ErrorContext(ctx, "failed to create tenant", "err", err)
		return nil, errorx.Wrap(err, errorx.ErrCodeInternalServer, "failed to create tenant")
	}
	key, err := tr.generateApiKey(tenant.ID)
	if err != nil {
		return nil, errorx.Wrap(err, errorx.ErrCodeInternalServer, "failed to generate api key")
	}
	tenant.ApiKey = key
	if err := tr.tenantRepo.Update(ctx, tenant); err != nil {
		slog.ErrorContext(ctx, "failed to update tenant", "err", err)
		return nil, errorx.Wrap(err, errorx.ErrCodeInternalServer, "failed to update tenant")
	}
	return tenant, nil
}

// generateApiKey generates a secure API key for tenant authentication
func (tr *TenantService) generateApiKey(tenantID uint) (string, error) {
	// 1. Convert tenant_id to bytes
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(tenantID))

	// 2. Encrypt tenant_id using AES-GCM
	block, err := aes.NewCipher(apiKeySecret())
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, idBytes, nil)

	// 3. Combine nonce and ciphertext, then encode with base64
	combined := append(nonce, ciphertext...)
	encoded := base64.RawURLEncoding.EncodeToString(combined)

	// Create final API Key in format: sk-{encrypted_part}
	return "sk-" + encoded, nil
}
