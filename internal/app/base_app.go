package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"silo/internal/config"
	"silo/pkg/echo_handle"
	"silo/pkg/echo_handle/middleware"
	"silo/pkg/logger"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"gorm.io/gorm"

	// https://www.elastic.co/guide/en/apm/agent/go/current/builtin-modules.html#builtin-modules-apmecho
	"go.elastic.co/apm/module/apmechov4/v2"
)

type BaseApp interface {
	SetServerAddress(address string)
	Run() error
}

// BaseAppImpl 基础应用结构，包含通用功能
type BaseAppImpl struct {
	e              *echo.Echo
	serverAddress  string
	shutdownCtx    context.Context
	shutdownCancel context.CancelFunc
	wg             sync.WaitGroup
	db             *gorm.DB
}

// NewBaseApp 创建基础应用
func NewBaseApp(e *echo.Echo, db *gorm.DB, cfg *config.CommonConfig) BaseApp {
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	baseApp := &BaseAppImpl{
		e:              e,
		shutdownCtx:    shutdownCtx,
		shutdownCancel: shutdownCancel,
		db:             db,
	}
	baseApp.initCommonEcho(cfg)
	return baseApp
}

// SetServerAddress 设置服务器地址
func (a *BaseAppImpl) SetServerAddress(address string) {
	a.serverAddress = address
}

// Run 运行应用，包含优雅关闭逻辑
func (a *BaseAppImpl) Run() error {
	// 设置优雅关闭
	a.setupGracefulShutdown()

	// 启动服务器
	serverErr := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				serverErr <- fmt.Errorf("server panic: %v", r)
			}
		}()
		serverErr <- a.e.Start(a.getServerAddress())
	}()

	// 等待服务器启动错误或关闭信号
	select {
	case err := <-serverErr:
		return err // 服务器启动失败或运行时错误
	case <-a.shutdownCtx.Done():
		// 收到关闭信号，执行优雅关闭
		a.gracefulShutdown()
		return nil // 正常关闭
	}
}

// InitCommonEcho 初始化通用Echo配置
func (a *BaseAppImpl) initCommonEcho(cfg *config.CommonConfig) {
	e := a.e
	e.HideBanner = true

	// 设置错误处理器
	e.HTTPErrorHandler = echo_handle.ErrorHandler

	// logger
	e.Logger = logger.NewEchoLogger(slog.Default())

	// validator
	e.Validator = echo_handle.NewValidator()

	// Remove trailing slash middleware removes a trailing slash from the request URI
	e.Pre(echoMiddleware.RemoveTrailingSlash())

	// 默认使用Recover中间件
	if cfg.Apm.Enabled {
		// recover panic
		e.Use(apmechov4.Middleware())
	} else {
		e.Use(echoMiddleware.Recover())
	}

	// RateLimiter
	e.Use(echoMiddleware.RateLimiter(
		echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(cfg.Api.BaseApiConfig.RateLimit.Rate)),
	))

	// CORS
	if cfg.Api.BaseApiConfig.CORS {
		//default cors
		e.Use(echoMiddleware.CORS())
	}

	// enable request id
	e.Use(echoMiddleware.RequestID())
	// logger go first
	e.Use(middleware.NewBodyDump())
}

// setupGracefulShutdown 设置优雅关闭
func (a *BaseAppImpl) setupGracefulShutdown() {
	// 监听信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			slog.Info("Received signal, starting graceful shutdown", "signal", sig.String())
			a.shutdownCancel()
		}
	}()
}

// GracefulShutdown 优雅关闭服务器
func (a *BaseAppImpl) gracefulShutdown() {
	// 设置关闭超时时间
	shutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	slog.Info("Shutting down server...")

	// 关闭HTTP服务器，等待现有请求完成
	if err := a.e.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("Server shutdown completed")
	}

	// 等待所有goroutine完成
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("All goroutines completed")
	case <-ctx.Done():
		slog.Warn("Timeout waiting for goroutines to complete")
	}
}

// getServerAddress 获取服务器地址
func (a *BaseAppImpl) getServerAddress() string {
	if a.serverAddress == "" {
		return ":8080" // 默认地址
	}
	return a.serverAddress
}
