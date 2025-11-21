/*
Copyright © 2025 lixw
*/
package v1

// TenantRequest tenant request
type TenantRequest struct {
	// tenant name
	Name string `json:"name"`
	// tenant description
	Description string `json:"description"`
}
