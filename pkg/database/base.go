/*
Copyright © 2025 lixw
*/
package database

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Status    uint8          `json:"status" gorm:"default:1;comment:状态:0=禁用,1=启用"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:删除时间" swaggerignore:"true"`
}
