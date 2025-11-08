package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// DefaultConfigMap 默认配置映射
func DefaultConfigMap() map[string]interface{} {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	serverID := uuid.New().String()
	logDir := dir + "/logs/"

	defaultsMap := map[string]interface{}{
		// app 默认配置
		"app.name":        "silo-website",
		"app.environment": "development",
		"app.debug":       true,
		"app.timezone":    "Asia/Singapore",
		"app.server_id":   serverID,
		"app.server_host": "0.0.0.0",
		"app.server_port": 8080,
		"app.language":    "en-US",

		// platform 默认配置
		"platform.addr":            "http://localhost:8080",
		"platform.timeout":         "30s",
		"platform.retry_count":     3,
		"platform.retry_wait_time": "1s",
		"platform.api_key":         "",
		"platform.auth_token":      "",
		"platform.debug":           true,

		// cache 默认配置
		"cache.mutex_lock_ttl":  10,
		"cache.cache_ttl":       3600,
		"cache.cache_partition": 10,

		// postgres 默认配置
		"postgres.debug":            true,
		"postgres.dsn":              "postgres://silo:silo@localhost:15432/silo?sslmode=disable&TimeZone=Asia/Singapore",
		"postgres.prepare_stmt":     true,
		"postgres.max_idle_conns":   10,
		"postgres.max_open_conns":   100,
		"postgres.monitor_interval": "60s",

		// redis 默认配置
		"redis.cluster_enabled": false,
		"redis.addr":            "localhost:6379",
		"redis.password":        "",

		// log 默认配置
		"log.level":                   "debug",
		"log.struct_log":              true,
		"log.pretty":                  true,
		"log.no_color":                false,
		"log.with_server_info":        true,
		"log.console_logging_enabled": true,
		"log.file_logging_enabled":    true,
		"log.directory":               logDir,
		"log.filename":                "app.log",
		"log.file_permission":         "644",
		"log.max_size":                100,
		"log.max_backups":             3,
		"log.max_age":                 28,
		"log.time_field_format":       "2006-01-02 15:04:05.000",
		"log.timestamp_field_name":    "timestamp",
		"log.level_field_name":        "level",
		"log.message_field_name":      "message",
		"log.error_field_name":        "error",

		// api 默认配置
		"api.base_api_config.skip_auth":           false,
		"api.base_api_config.swagger":             true,
		"api.base_api_config.port":                8080,
		"api.base_api_config.cors":                true,
		"api.base_api_config.request_log_enabled": true,
		"api.base_api_config.rate_limit.enabled":  true,
		"api.base_api_config.rate_limit.rate":     100,
		"api.base_api_config.front_domain":        "",
		"api.base_api_config.request_ttl":         "60s",
		"api.base_api_config.request_key":         "",

		// admin 特有配置
		"api.auth.signing_key":     "your-secret-key",
		"api.auth.online_time":     7200,
		"api.super_admin.name":     "admin",
		"api.super_admin.email":    "admin@example.com",
		"api.super_admin.password": "admin123",

		// apm 默认配置
		"apm.enabled":      false,
		"apm.server_url":   "",
		"apm.secret_token": "",
		"apm.service_name": "",
	}

	return defaultsMap
}

// NewConfigWithPrefix 创建配置，使用默认前缀 SILO，环境变量优先级最高
func NewConfigWithPrefix(configFile string, prefix string) *viper.Viper {
	v := viper.New()

	// 首先设置默认值
	defaults := DefaultConfigMap()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	// 然后设置环境变量
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // 自动读取环境变量

	// 最后读取配置文件（如果存在）
	if configFile != "" && checkFileExists(configFile) {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			log.Printf("Error reading config file, %s", err)
		}
	}

	return v
}

// NewConfig 创建配置，环境变量优先级最高
func NewConfig(envPrefix string, configFile string) *viper.Viper {
	v := viper.New()

	// 首先设置默认值
	defaults := DefaultConfigMap()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	// 然后设置环境变量前缀和键替换规则
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // 自动读取环境变量

	// 最后读取配置文件（如果存在）
	if configFile != "" && checkFileExists(configFile) {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			log.Printf("Error reading config file, %s", err)
		}
	}

	return v
}

// checkFileExists 检查文件是否存在
func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
