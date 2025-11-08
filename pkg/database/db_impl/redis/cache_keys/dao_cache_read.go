package cache_keys

const (
	// ModelReadKeyF 开启CacheReadOnly的Model的缓存Key, 格式：{model_read:{redis分片ID}}:{model cache key}
	ModelReadKeyF = "model_cache_read:%s"
)
