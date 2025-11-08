package config

// APMConfig .
type APMConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServerURL   string `mapstructure:"server_url"`
	SecretToken string `mapstructure:"secret_token"`
	ServiceName string `mapstructure:"service_name"`
}
