package main

import (
	"flag"
	"log"
	"silo/internal/config"
	"silo/internal/infra/model"
	"silo/pkg/database/db_impl/postgres"
	"silo/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func main() {
	// 解析命令行参数
	needInit := flag.Bool("init", false, "Initialize database")
	configPath := flag.String("config", "../../configs/website.yaml", "Path to config file")
	prefix := flag.String("p", "QQFUN", "Environment variable prefix (default: QQFUN)")
	flag.Parse()

	// 初始化配置
	cfg := initConfig(*configPath, *prefix)

	// 初始化日志
	initLog(cfg)

	// 初始化数据库
	if *needInit {
		db := postgres.NewPostgres(cfg.Postgres)
		if err := initDB(db); err != nil {
			panic(err)
		}
		return
	}

	// 初始化应用
	e := echo.New()
	app, err := InitApp(&cfg, e)
	if err != nil {
		panic(err)
	}

	// 运行应用
	if err = app.Run(); err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}
}

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

// initDB 初始化数据库
func initDB(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	return model.Migration(db)
}
