package cache_keys

const (
	// ModelKeyF 开启CacheAble的Model的缓存Key, 格式：{model:{redis分片ID}}:{model cache key}
	ModelKeyF = "{model:%d}:%s"
	// ModelMutexKeyF Model互斥锁Key, 格式：{model:{redis分片ID}}:mutex:{model cache key}
	ModelMutexKeyF = "{model:%d}:mutex:%s"
)
