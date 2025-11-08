package model

import "silo/pkg/types/isotime"

const TableNameSysMigration = "sys_migration"

//go:generate go run accessor SysMigration -f sys_migration
type SysMigration struct {
	Name      string          `json:"name" gorm:"column:name;primary_key"`
	CreatedAt isotime.ISOTime `json:"created_at"`
}

// TableName overrides table name
func (obj *SysMigration) TableName() string {
	return TableNameSysMigration
}

// IDMap for filter
func (obj *SysMigration) IDMap() map[string]any {
	return map[string]any{
		"name": obj.Name,
	}
}
