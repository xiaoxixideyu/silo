package postgres

import (
	"fmt"
	"log/slog"
	"silo/pkg/database/config"
	"silo/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	postgres "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
)

// NewPostgres .
func NewPostgres(cfg config.PostgresConfig) *gorm.DB {
	var db *gorm.DB
	var err error
	dsn := cfg.DSN
	gc := &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
		PrepareStmt:              cfg.PrepareStmt,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
	gc.Logger = logger.NewGormLogger(slog.Default().With("module", "gorm"), cfg.Debug)

	db, err = gorm.Open(postgres.Open(dsn), gc)
	if err != nil {
		msg := fmt.Sprintf("[postgres] connect db error: %s", err.Error())
		slog.Error(msg)
		panic(msg)
	}
	sqlDB, err := db.DB()
	if err != nil {
		msg := fmt.Sprintf("[postgres]failed to get db error: %s", err.Error())
		slog.Error(msg)
		panic(msg)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	return db
}
