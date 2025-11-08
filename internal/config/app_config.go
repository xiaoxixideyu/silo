package config

import "time"

// AppConfig 应用基础配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
	Timezone    string `mapstructure:"timezone"`
	ServerId    string `mapstructure:"server_id"`
	ServerHost  string `mapstructure:"server_host"`
	ServerPort  int    `mapstructure:"server_port"`
	Language    string `mapstructure:"language"`
}

// BaseAPIConfig 基础API配置
type BaseAPIConfig struct {
	SkipAuth          bool `mapstructure:"skip_auth"`
	Swagger           bool
	Port              int
	CORS              bool
	RequestLogEnabled bool `mapstructure:"request_log_enabled"`
	RateLimit         struct {
		Enabled bool
		// {rate number} requests/sec
		Rate int
	} `mapstructure:"rate_limit"`
	FrontDomain string        `mapstructure:"front_domain"`
	RequestTTL  time.Duration `mapstructure:"request_ttl"`
	RequestKey  string        `mapstructure:"request_key"`
}
