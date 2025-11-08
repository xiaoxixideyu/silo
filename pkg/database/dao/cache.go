package dao

import (
	"context"
	"encoding/json"
	"log/slog"
	"silo/pkg/database/config"
	"silo/pkg/database/db_impl/redis/cache_keys"
	"silo/pkg/database/db_impl/redis/redlock"
	"silo/pkg/database/interfaces"
	"silo/pkg/database/transaction"
	"silo/pkg/utils"
	"time"

	"github.com/avast/retry-go"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

/**
* Cache-Aside caching mechanism
* 1. Read operations (FindCache):
   - Read cache first, return directly if cache exists
   - If cache doesn't exist, use singleflight to prevent cache breakdown
   - Read data from database and write to cache
* 2. Write operations:
   - Create operations (FindOrCreate):
     * Check cache first, return directly if cache hit
     * If cache miss, find or create from database
     * If data is found, write to cache immediately
     * If data is newly created and in transaction, write to cache after transaction commits (to prevent transaction rollback)
     * If data is newly created and not in transaction, write to cache immediately
   - Update operations (SaveCache):
     * Write to database first
     * After successful write, delete corresponding cache (to ensure data consistency)
     * If in transaction, delete cache after transaction commits
**/

// CacheNotFoundError cache not found error
const CacheNotFoundError = redis.Nil

// NewDaoCache .
func NewDaoCache(db *gorm.DB, cacheClient interfaces.CacheClient, config config.ModelCacheConfig) interfaces.DaoCache {
	return &DAOCacheImpl{
		DAOImpl:     DAOImpl{db: db},
		cacheClient: cacheClient,
		config:      config,
		redlock:     redlock.NewRedlock(cacheClient),
	}
}

// DAOCacheImpl .
type DAOCacheImpl struct {
	DAOImpl
	cacheClient interfaces.CacheClient
	config      config.ModelCacheConfig
	sf          singleflight.Group
	redlock     redlock.Redlock
}

// FindCache implements Cache-Aside pattern read operation
func (c *DAOCacheImpl) FindCache(ctx context.Context, model interfaces.ModelWithCache, out interfaces.ModelWithCache) error {
	cacheKey := model.CacheKey()

	// Use singleflight to prevent cache breakdown
	_, err, _ := c.sf.Do(cacheKey, func() (any, error) {
		// 1. Try to read from cache first
		exist, err := c.findCache(ctx, cacheKey, out)
		if err != nil {
			slog.Warn("[DAOCacheImpl] findCache error", "cacheKey", cacheKey, "error", err)
			return nil, err
		}

		// 2. Cache exists, return directly
		if exist {
			return nil, nil
		}

		// 3. Cache doesn't exist, query from database and set cache
		if err = c.findFromDBAndSetCache(ctx, model, out, cacheKey, true); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// FindOrCreateCache implements Cache-Aside pattern find or create operation
func (c *DAOCacheImpl) FindOrCreateCache(ctx context.Context, model interfaces.ModelWithCache, out interfaces.ModelWithCache) (bool, error) {
	cacheKey := model.CacheKey()

	// Use singleflight to prevent concurrent database queries and duplicate creation within single process
	val, err, _ := c.sf.Do(cacheKey, func() (any, error) {
		// 1. Try to read from cache first
		exist, err := c.findCache(ctx, cacheKey, out)
		if err != nil {
			slog.Warn("[DAOCacheImpl] findCache(findOrCreate) error", "model", model.TableName(), "id", model.IDMap(), "error", err)
			return false, err
		}

		// 2. Cache exists, return directly
		if exist {
			return false, nil
		}

		// 3. Check if in transaction
		tx := transaction.FromContext(ctx)
		if tx != nil {
			// In transaction: no distributed lock, query database directly
			return c.findOrCreateFromDB(ctx, model, out)
		}

		// 4. Not in transaction: use redlock
		mutexKey := c.getMutexKey(cacheKey)
		mutex := c.redlock.NewMutex(mutexKey, redlock.WithExpiry(time.Duration(c.config.MutexLockTTL)*time.Second))
		// Use LockFunc for auto-renewal, return error if lock acquisition fails
		var created bool
		if err = mutex.LockFunc(ctx, func(ctx context.Context) error {
			// 5. Check cache again after acquiring lock (double check)
			exist, err = c.findCache(ctx, cacheKey, out)
			if err != nil {
				slog.Warn("[DAOCacheImpl] findCache(findOrCreate) after lock error", "cacheKey", cacheKey, "error", err)
				return err
			}
			if exist {
				created = false
				return nil
			}

			// 6. Find or create from database
			created, err = c.findOrCreateFromDB(ctx, model, out)
			return err
		}); err != nil {
			slog.Error("[DAOCacheImpl] redlock(findOrCreate) error", "cacheKey", cacheKey, "error", err)
			return false, err
		}

		return created, nil
	})

	if err != nil {
		return false, err
	}

	return val.(bool), nil
}

// SaveCache implements Cache-Aside pattern write operation
func (c *DAOCacheImpl) SaveCache(ctx context.Context, model interfaces.ModelWithCache) error {
	changes := model.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	// Check if in transaction
	tx := transaction.FromContext(ctx)
	if tx == nil {
		// Non-transaction mode: write to database directly, then delete cache
		if err := c.Save(ctx, model); err != nil {
			return err
		}

		// Delete cache (with retry)
		if err := c.deleteCacheWithRetry(ctx, model.CacheKey()); err != nil {
			slog.Error("[DAOCacheImpl] deleteCache error", "model", model.TableName(), "id", model.IDMap(), "error", err)
			return err
		}

		return nil
	}

	// Transaction mode: register hook to delete cache after transaction commits
	if err := c.Save(ctx, model); err != nil {
		return err
	}

	tx.AfterCommit(func() {
		bgCtx := context.Background()
		if err := c.deleteCacheWithRetry(bgCtx, model.CacheKey()); err != nil {
			slog.Error("[DAOCacheImpl] deleteCache error after commit", "model", model.TableName(), "id", model.IDMap(), "error", err)
		}
	})

	return nil
}

// ResetCacheTTL reset cache expiration time
func (c *DAOCacheImpl) ResetCacheTTL(ctx context.Context, obj interfaces.ModelWithCache) error {
	modelKey := c.getModelKey(obj.CacheKey())

	// Use Redis Expire command directly to set TTL
	err := c.cacheClient.Expire(ctx, modelKey, time.Duration(c.config.CacheTTL)*time.Second).Err()
	if err != nil {
		slog.Error("[DAOCacheImpl] ResetCacheTTL error", "error", err)
		return err
	}

	return nil
}

// findCache get data from cache
func (c *DAOCacheImpl) findCache(ctx context.Context, cacheKey string, out interfaces.ModelWithCache) (bool, error) {
	modelKey := c.getModelKey(cacheKey)

	// Use Redis Get command directly to get value
	value, err := c.cacheClient.Get(ctx, modelKey).Result()
	if err != nil {
		if c.cacheClient.IsNil(err) {
			return false, nil
		}
		return false, err
	}

	if value == "" {
		return false, nil
	}

	if err := json.Unmarshal([]byte(value), out); err != nil {
		return false, err
	}

	return true, nil
}

// setCache set cache data
func (c *DAOCacheImpl) setCache(ctx context.Context, model interfaces.ModelWithCache, ttl int64) error {
	v, err := json.Marshal(model)
	if err != nil {
		return err
	}

	modelKey := c.getModelKey(model.CacheKey())
	err = c.cacheClient.Set(ctx, modelKey, string(v), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

// deleteCache delete cache data
func (c *DAOCacheImpl) deleteCache(ctx context.Context, cacheKey string) error {
	modelKey := c.getModelKey(cacheKey)
	return c.cacheClient.Del(ctx, modelKey).Err()
}

// deleteCacheWithRetry delete cache (with retry mechanism)
func (c *DAOCacheImpl) deleteCacheWithRetry(ctx context.Context, cacheKey string) error {
	return retry.Do(
		func() error {
			return c.deleteCache(ctx, cacheKey)
		},
		retry.LastErrorOnly(true),
		retry.RetryIf(func(e error) bool {
			return !c.cacheClient.IsNil(e)
		}),
		retry.Attempts(3),
		retry.Delay(100*time.Millisecond),
		retry.OnRetry(func(n uint, err error) {
			slog.Warn("[DAOCacheImpl] delete cache retry", "cacheKey", cacheKey, "attempt", n, "error", err.Error())
		}),
	)
}

// findFromDBAndSetCache query data from database and set cache
func (c *DAOCacheImpl) findFromDBAndSetCache(ctx context.Context, model interfaces.ModelWithCache, out interfaces.ModelWithCache, cacheKey string, setCache bool) error {
	// Query data from database
	if err := c.FindOne(ctx, model.TableName(), NewRequest().FilterM(model.IDMap()), out); err != nil {
		return ErrorMap(err, model.TableName())
	}

	// Set cache if needed
	if setCache {
		if err := c.setCache(ctx, out, c.config.CacheTTL); err != nil {
			slog.Warn("[DAOCacheImpl] setCache error", "error", err)
		}
	}

	return nil
}

// findOrCreateFromDB find or create from database and handle cache
func (c *DAOCacheImpl) findOrCreateFromDB(ctx context.Context, model interfaces.ModelWithCache, out interfaces.ModelWithCache) (bool, error) {
	create, err := c.FindOrCreate(ctx, model.TableName(), model, out)
	if err != nil {
		return false, err
	}

	// Handle cache write
	tx := transaction.FromContext(ctx)
	if tx != nil && create {
		// Created in transaction: delay cache write until after transaction commits
		tx.AfterCommit(func() {
			bgCtx := context.Background()
			if err := c.setCache(bgCtx, out, c.config.CacheTTL); err != nil {
				slog.Warn("[DAOCacheImpl] setCache(findOrCreate) after commit error", "model", model.TableName(), "id", model.IDMap(), "error", err)
			}
		})
	} else {
		// Created outside transaction or found existing data: write to cache immediately
		if err := c.setCache(ctx, out, c.config.CacheTTL); err != nil {
			slog.Warn("[DAOCacheImpl] setCache(findOrCreate) error", "model", model.TableName(), "id", model.IDMap(), "error", err)
		}
	}

	return create, nil
}

// -- Redis cluster compatible mode --

// getPartition for Redis cluster mode, distribute cache data to different Redis instances
func (c *DAOCacheImpl) getPartition(cacheKey string) int {
	return utils.HashToPartition(cacheKey, c.config.CachePartition)
}

// getModelKey get model cache key
func (c *DAOCacheImpl) getModelKey(cacheKey string) string {
	partition := c.getPartition(cacheKey)
	return c.cacheClient.Keyf(cache_keys.ModelKeyF, partition, cacheKey)
}

// getMutexKey get mutex lock key
func (c *DAOCacheImpl) getMutexKey(cacheKey string) string {
	partition := c.getPartition(cacheKey)
	return c.cacheClient.Keyf(cache_keys.ModelMutexKeyF, partition, cacheKey)
}
