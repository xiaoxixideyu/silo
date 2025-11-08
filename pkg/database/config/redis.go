package config

// RedisConfig for redis
type RedisConfig struct {
	ClusterEnabled bool   `mapstructure:"cluster_enabled"`
	Addr           string `mapstructure:"addr"`
	Password       string `mapstructure:"password"`
	Prefix         string `mapstructure:"prefix"`
}
