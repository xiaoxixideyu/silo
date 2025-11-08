package logger

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

func NewGormLogger(l *slog.Logger, debug bool) logger.Interface {
	return &GormLogger{
		logger: l,
		debug:  debug,
	}
}

type GormLogger struct {
	logger *slog.Logger
	debug  bool
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.debug {
		sql, rowsAffected := fc()
		l.logger.Debug(fmt.Sprintf("sql: %s, rowsAffected: %d, err: %v", sql, rowsAffected, err))
	}
}
