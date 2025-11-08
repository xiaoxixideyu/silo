package config

import "time"

// PlatformConfig platform configuration
type PlatformConfig struct {
	// Basic configuration
	Addr string `mapstructure:"addr"` // Platform service address
	// Timeout configuration
	Timeout time.Duration `mapstructure:"timeout"` // Request timeout, default 30 seconds
	// Retry configuration
	RetryCount    int           `mapstructure:"retry_count"`     // Retry count, default 3 times
	RetryWaitTime time.Duration `mapstructure:"retry_wait_time"` // Retry wait time, default 1 second
	// Authentication configuration
	APIKey    string `mapstructure:"api_key"`    // API key
	AuthToken string `mapstructure:"auth_token"` // Authentication token
	// Debug configuration
	Debug bool `mapstructure:"debug"` // Whether to enable debug mode
}

// GetDefaultConfig get default configuration
func GetDefaultConfig() PlatformConfig {
	return PlatformConfig{
		Addr:          "http://localhost:8081",
		Timeout:       30 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
		Debug:         false,
	}
}
