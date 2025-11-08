package config

import "time"

type PostgresConfig struct {
	Debug bool   `mapstructure:"debug"`
	DSN   string `mapstructure:"dsn"`
	// PrepareStmt enables prepared statements.
	// turn to false, if pgsql return "cached plan must not change result type (SQLSTATE 0A000)""
	PrepareStmt bool `mapstructure:"prepare_stmt"`
	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
	//
	// If n <= 0, no idle connections are retained.
	//
	// The default max idle connections is currently 2. This may change in
	// a future release.
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MonitorInterval time.Duration `mapstructure:"monitor_interval"` // 监控间隔(秒)
}
