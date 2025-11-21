/*
Copyright © 2025 lixw
*/
package model

import (
	"github.com/ethanli-dev/go-app-layout/pkg/database"
)

type Tenant struct {
	database.Model
	Name        string `json:"name" gorm:"column:name;size:127;not null;comment:租户名称"`
	Description string `json:"description" gorm:"column:description;size:511;comment:租户描述"`
	ApiKey      string `json:"api_key" gorm:"column:api_key;size:255;comment:API密钥"`
}

func (*Tenant) TableName() string {
	return "tenant"
}
