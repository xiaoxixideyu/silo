package config

import (
	databaseConfig "silo/pkg/database/config"
	platformConfig "silo/pkg/platform/config"
)

// Config 配置接口
type Config interface {
	GetPostgresConfig() databaseConfig.PostgresConfig
	GetRedisConfig() databaseConfig.RedisConfig
	GetModelCacheConfig() databaseConfig.ModelCacheConfig
	GetPlatformConfig() platformConfig.PlatformConfig
}
