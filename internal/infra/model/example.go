package model

import (
	"fmt"
	"silo/pkg/database/seedwork"
)

const TableNameExample = "example"

//go:generate go run accessor Example -f example
type Example struct {
	seedwork.ModelWithCache
	seedwork.WithTimestamp
	seedwork.WithSoftDelete
	ID   int64  `json:"id" gorm:"primary_key"`
	Info string `json:"info"`
	Val  int64  `json:"val"`
}

// TableName overrides table name
func (obj *Example) TableName() string {
	return TableNameExample
}

// IDMap for filter
func (obj *Example) IDMap() map[string]any {
	return map[string]any{
		"id": obj.ID,
	}
}

// CacheKey .
func (obj *Example) CacheKey() string {
	return fmt.Sprintf("%s:%d", obj.TableName(), obj.ID)
}
