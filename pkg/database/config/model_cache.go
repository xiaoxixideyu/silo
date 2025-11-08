package config

// ModelCacheConfig for model cache
type ModelCacheConfig struct {
	MutexLockTTL   int64 `mapstructure:"mutex_lock_ttl"`  // 互斥锁过期时间
	CacheTTL       int64 `mapstructure:"cache_ttl"`       // 缓存过期时间
	CachePartition int   `mapstructure:"cache_partition"` // 缓存分区数，此值不能随意修改会导致缓存保存到不同的节点中，导致缓存数据不一致
}
