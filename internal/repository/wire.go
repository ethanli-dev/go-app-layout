/*
Copyright © 2025 lixw
*/
package repository

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTenantRepository)
