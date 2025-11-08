package interfaces

// Model .
type Model interface {
	// DB table name
	TableName() string
	// IDMap for filter (used by findOrCreate and delete)
	IDMap() map[string]any
	// Changes
	GetChanges() map[string]any
	// SetChanges
	Update(name string, value any)
}

// ModelWithCache .
type ModelWithCache interface {
	Model
	WithCacheable
}

// ModelWithCacheRead .
type ModelWithCacheRead interface {
	Model
	CacheKey() string
}

type WithVersion interface {
	IncVersion()
	GetVersion() int64
}

// WithCacheable .
type WithCacheable interface {
	WithVersion
	CacheKey() string
	SetVersion(v int64)
}
