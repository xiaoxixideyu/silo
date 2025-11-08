package interfaces

import (
	"context"
)

//go:generate mockgen -source=dao_cache.go -destination=mocks/dao_cache.go

// DaoCache .
type DaoCache interface {
	Dao
	FindCache(ctx context.Context, model ModelWithCache, out ModelWithCache) error
	FindOrCreateCache(ctx context.Context, model ModelWithCache, out ModelWithCache) (create bool, err error)
	// SaveCache DAOCache数据保存，实现Cache-Aside模式（写操作删缓存）
	SaveCache(ctx context.Context, model ModelWithCache) error
	// ResetCacheTTL 重置缓存过期时间
	ResetCacheTTL(ctx context.Context, model ModelWithCache) error
}
