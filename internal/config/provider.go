package config

import (
	databaseConfig "silo/pkg/database/config"
	platformConfig "silo/pkg/platform/config"

	"github.com/google/wire"
)

// ConfigSet 包含所有配置相关的wire set
var ConfigSet = wire.NewSet(
	ProvidePostgresConfig,
	ProvideRedisConfig,
	ProvideModelCacheConfig,
	ProvidePlatformConfig,
)

func ProvidePostgresConfig(cfg Config) databaseConfig.PostgresConfig {
	return cfg.GetPostgresConfig()
}

func ProvideRedisConfig(cfg Config) databaseConfig.RedisConfig {
	return cfg.GetRedisConfig()
}

func ProvideModelCacheConfig(cfg Config) databaseConfig.ModelCacheConfig {
	return cfg.GetModelCacheConfig()
}

func ProvidePlatformConfig(cfg Config) platformConfig.PlatformConfig {
	return cfg.GetPlatformConfig()
}
