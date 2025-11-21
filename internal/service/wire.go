/*
Copyright © 2025 lixw
*/
package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTenantService)
