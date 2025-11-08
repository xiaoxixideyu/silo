package config

// LogConfig LogConfig
type LogConfig struct {
	Level string `mapstructure:"level"`
	// whether to log in structured format or not.
	StructLog bool `mapstructure:"struct_log"`
	// Formats JSON with indentation for better readability, works if StructLog is false
	Pretty        bool `mapstructure:"pretty"`
	NoColor       bool `mapstructure:"no_color"`
	WithServeInfo bool `mapstructure:"with_server_info"`
	// Enable console logging
	ConsoleLoggingEnabled bool `mapstructure:"console_logging_enabled"`
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `mapstructure:"file_logging_enabled"`
	// Directory to log to to when filelogging is enabled
	Directory string `mapstructure:"directory"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `mapstructure:"filename"`
	// File permission 644 for rw-r--r--
	FilePermission string `mapstructure:"file_permission"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `mapstructure:"max_size"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `mapstructure:"max_age"`
	// Log format
	TimeFieldFormat    string `mapstructure:"time_Field_format"`
	TimestampFieldName string `mapstructure:"timestamp_field_name"`
	LevelFieldName     string `mapstructure:"level_field_name"`
	MessageFieldName   string `mapstructure:"message_field_name"`
	ErrorFieldName     string `mapstructure:"error_field_name"`
	ServerType         string
	ServerID           string
}
