package seedwork

import (
	"silo/pkg/database/interfaces"
	"silo/pkg/types/isotime"
)

var _ (interfaces.WithTimestamp) = (*WithTimestamp)(nil)

const (
	TableColumnID        = "id"
	TableColumnCreatedAt = "created_at"
	TableColumnDeleteAt  = "deleted_at"
	TableColumnUpdatedAt = "updated_at"
	TableColumnVersion   = "version"
)

// WithTimestamp with createdAt and updatedAt
type WithTimestamp struct {
	CreatedAt isotime.ISOTime `json:"created_at,omitempty"`
	UpdatedAt isotime.ISOTime `json:"updated_at,omitempty"`
}

// GetUpdatedAt .
func (w *WithTimestamp) GetUpdatedAt() isotime.ISOTime {
	return w.UpdatedAt
}

// SetUpdatedAt .
func (w *WithTimestamp) WithTimestampSetUpdatedAt() {
	w.UpdatedAt = isotime.Now()
}

// GetCreatedAt .
func (w *WithTimestamp) GetCreatedAt() isotime.ISOTime {
	return w.CreatedAt
}

// SetCreatedAt .
func (w *WithTimestamp) WithTimestampSetCreatedAt() {
	w.CreatedAt = isotime.Now()
}

/************ WithSoftDelete ************/

// WithSoftDelete with soft delete enabled
type WithSoftDelete struct {
	DeletedAt isotime.ISOTime `json:"deleted_at,omitempty" bson:"deletedAt,omitempty"`
}

// SetDeletedAt .
func (w *WithSoftDelete) SetDeletedAt() {
	w.DeletedAt = isotime.Now()
}

func (w *WithSoftDelete) IsDeleted() bool {
	return !w.DeletedAt.IsZero()
}

/************ WithVersion ************/

// WithVersion .
type WithVersion struct {
	Version int64 `json:"version,omitempty" gorm:"column:version;default:1"`
}

// IncVersion .
func (w *WithVersion) IncVersion() {
	w.Version++
}

// GetVersion .
func (w *WithVersion) GetVersion() int64 {
	return w.Version
}

/************ ModelWithCache ************/

// ModelWithCache .
type ModelWithCache struct {
	Model       `json:"-"`
	WithVersion `accessor:",inline"`
}

// CacheKey .
func (b *ModelWithCache) CacheKey() string {
	panic("not implement")
}

// SetVersion SetVersion
func (b *ModelWithCache) SetVersion(v int64) {
	b.Version = v
	b.Update("version", v)
}

// // IncVersion IncVersion
// func (b *ModelWithCache) IncVersion() int64 {
// 	b.Version++
// 	b.Update("version", b.Version)

// 	return b.Version
// }

/************ Model ************/

// Model .
type Model struct {
	changes map[string]any
}

// GetChanges .
func (obj *Model) GetChanges() map[string]any {
	if obj.changes == nil {
		return nil
	}
	result := obj.changes
	obj.changes = nil
	return result
}

// Update .
func (obj *Model) Update(name string, value any) {
	if obj.changes == nil {
		obj.changes = make(map[string]any)
	}
	obj.changes[name] = value
}

/************ ModelWithCacheRead ************/

// // ModelWithCacheRead .
// type ModelWithCacheRead struct {
// 	Model `json:"-"`
// }

// // CacheKey .
// func (b *ModelWithCacheRead) CacheKey() string {
// 	panic("not implement")
// }
