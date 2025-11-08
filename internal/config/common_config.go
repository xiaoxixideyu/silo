package config

import (
	database "silo/pkg/database/config"
	loggerConfig "silo/pkg/logger/config"
	platformConfig "silo/pkg/platform/config"
)

// CommonConfig 统一配置结构体，适用于admin和website
type CommonConfig struct {
	App      AppConfig                     `mapstructure:"app"`
	Platform platformConfig.PlatformConfig `mapstructure:"platform"`
	Cache    database.ModelCacheConfig     `mapstructure:"cache"`
	Postgres database.PostgresConfig       `mapstructure:"postgres"`
	Redis    database.RedisConfig          `mapstructure:"redis"`
	Apm      APMConfig                     `mapstructure:"apm"`
	Api      CommonAPIConfig               `mapstructure:"api"`
	Log      loggerConfig.LogConfig        `mapstructure:"log"`
}

// CommonAPIConfig 统一API配置，包含admin和website的所有配置
type CommonAPIConfig struct {
	BaseApiConfig BaseAPIConfig `mapstructure:"base_api_config"`
	Auth          *AuthConfig   `mapstructure:"auth,omitempty"`        // 仅admin使用
	SuperAdmin    *SuperAdmin   `mapstructure:"super_admin,omitempty"` // 仅admin使用
}

// AuthConfig admin认证配置
type AuthConfig struct {
	SigningKey string `mapstructure:"signing_key"`
	OnlineTime int    `mapstructure:"online_time"` // 默认在线时间，秒
}

// SuperAdmin 超级管理员配置
type SuperAdmin struct {
	Name     string `mapstructure:"name"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

// GetPostgresConfig 实现 DatabaseConfig 接口
func (c CommonConfig) GetPostgresConfig() database.PostgresConfig {
	return c.Postgres
}

// GetRedisConfig 实现 DatabaseConfig 接口
func (c CommonConfig) GetRedisConfig() database.RedisConfig {
	// 动态设置 Redis prefix 为 app name 和 app environment 的拼接
	c.Redis.Prefix = c.App.Name + ":" + c.App.Environment + ":"
	return c.Redis
}

// GetModelCacheConfig 实现 DatabaseConfig 接口
func (c CommonConfig) GetModelCacheConfig() database.ModelCacheConfig {
	return c.Cache
}

// GetPlatformConfig .
func (c CommonConfig) GetPlatformConfig() platformConfig.PlatformConfig {
	return c.Platform
}
