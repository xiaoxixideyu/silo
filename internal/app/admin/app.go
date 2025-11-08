package admin

import (
	"fmt"
	"silo/internal/app"
	"silo/internal/config"
	handler "silo/internal/interface/admin"

	_ "silo/api/admin" // This line is important for swagger

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// App admin应用结构
type App struct {
	app.BaseApp
	cfg          config.CommonConfig
	adminHandler handler.AdminHandler
}

// NewApp 创建admin应用
func NewApp(cfg *config.CommonConfig, e *echo.Echo, adminHandler handler.AdminHandler, db *gorm.DB) *App {
	baseApp := app.NewBaseApp(e, db, cfg)
	// 设置服务器地址
	baseApp.SetServerAddress(fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.App.ServerPort))
	// 初始化路由
	g := e.Group("")
	adminHandler.InitRouter(g)

	// 根据配置决定是否添加 Swagger 路由
	if cfg.Api.BaseApiConfig.Swagger {
		// 配置 Swagger 文档
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return &App{
		BaseApp:      baseApp,
		cfg:          *cfg,
		adminHandler: adminHandler,
	}
}
