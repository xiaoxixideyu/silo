package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// 解析命令行参数
	needInit := flag.Bool("init", false, "Initialize database")
	configPath := flag.String("config", "../../configs/admin.yaml", "Path to config file")
	prefix := flag.String("p", "SILO", "Environment variable prefix (default: SILO)")
	flag.Parse()

	// 初始化配置
	cfg := initConfig(*configPath, *prefix)

	// 初始化日志
	initLog(cfg)

	// 初始化数据库
	if *needInit {
		initAdmin(&cfg)
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
