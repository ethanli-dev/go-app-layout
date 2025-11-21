/*
Copyright © 2025 lixw
*/
package handler

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTenantHandler)
