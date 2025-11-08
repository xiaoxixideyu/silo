//go:build wireinject
// +build wireinject

package main

import (
	app "silo/internal/app/admin"
	"silo/internal/config"
	handler "silo/internal/interface/admin"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

func InitApp(cfg *config.CommonConfig, e *echo.Echo) (*app.App, error) {
	wire.Build(
		app.NewApp,
		handler.HandlerSet,
		config.ConfigSet,
		wire.Bind(new(config.Config), new(*config.CommonConfig)),
	)
	return &app.App{}, nil
}
