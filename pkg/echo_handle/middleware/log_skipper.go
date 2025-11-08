package middleware

import (
	"github.com/labstack/echo/v4"
)

// CtxKeyLogSkipper log skipper
const CtxKeyLogSkipper = "logSkipper"

// NewLogSkipperConfig .
func NewLogSkipperConfig() *LogSkipperConfig {
	return &LogSkipperConfig{}
}

// LogSkipperConfig .
type LogSkipperConfig struct {
	SkipFullLog    bool
	SkipRequestLog bool
	SkipFilterFunc func(c echo.Context) bool
}

// WithSkipFullLog .
func (c *LogSkipperConfig) WithSkipFullLog(v bool) *LogSkipperConfig {
	c.SkipFullLog = v
	return c
}

// WithSkipRequestLog .
func (c *LogSkipperConfig) WithSkipRequestLog(v bool) *LogSkipperConfig {
	c.SkipRequestLog = v
	return c
}

// WithSkipFilterFunc .
func (c *LogSkipperConfig) WithSkipFilterFunc(v func(c echo.Context) bool) *LogSkipperConfig {
	c.SkipFilterFunc = v
	return c
}

func (c *LogSkipperConfig) ShouldSkipFullLog(ctx echo.Context) bool {
	if c.SkipFilterFunc != nil {
		return c.SkipFilterFunc(ctx) && c.SkipFullLog
	}

	return c.SkipFullLog
}

func (c *LogSkipperConfig) ShouldSkipRequestLog(ctx echo.Context) bool {
	if c.SkipFilterFunc != nil {
		return c.SkipFilterFunc(ctx) && c.SkipRequestLog
	}

	return c.SkipRequestLog
}

// NewLogSkipper skip current api from logging
func NewLogSkipper(cfg ...*LogSkipperConfig) echo.MiddlewareFunc {
	var skipper *LogSkipperConfig

	if len(cfg) > 0 {
		skipper = cfg[0]
	} else {
		skipper = NewLogSkipperConfig()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(CtxKeyLogSkipper, skipper)

			return next(c)
		}
	}
}
