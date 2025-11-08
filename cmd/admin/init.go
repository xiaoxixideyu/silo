package main

import (
	"log"
	"silo/internal/config"
	"silo/internal/infra/model"
	"silo/pkg/database/db_impl/postgres"
	"silo/pkg/logger"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// initConfig 初始化配置
func initConfig(configPath string, prefix string) config.CommonConfig {
	// 使用新的配置初始化方式，环境变量优先级最高
	v := config.NewConfigWithPrefix(configPath, prefix)

	var cfg config.CommonConfig
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return cfg
}

// initLog 初始化日志
func initLog(cfg config.CommonConfig) {
	// 初始化默认日志器，使用配置文件中的日志配置
	logger.InitDefaultLogger(&cfg.Log)
}

// initAdmin 项目初始化
func initAdmin(cfg *config.CommonConfig) {
	db := postgres.NewPostgres(cfg.Postgres)
	if err := initDB(db, cfg); err != nil {
		panic(err)
	}
}

// initDB 初始化数据库
func initDB(db *gorm.DB, cfg *config.CommonConfig) error {
	if db == nil {
		return errors.New("db is nil")
	}

	// 执行数据库迁移
	if err := model.Migration(db); err != nil {
		return err
	}

	return nil
}
