//go:build wireinject
// +build wireinject

package main

import (
	app "silo/internal/app/website"
	"silo/internal/config"
	handler "silo/internal/interface/website"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

func InitApp(cfg *config.CommonConfig, e *echo.Echo) (*app.App, error) {
	wire.Build(
		app.NewApp,
		handler.HandlerSet,
		config.ConfigSet,
		// 将 config.CommonConfig 转换为 config.Config
		wire.Bind(new(config.Config), new(*config.CommonConfig)),
	)
	return &app.App{}, nil
}
