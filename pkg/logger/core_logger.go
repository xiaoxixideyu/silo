package logger

import (
	"fmt"
	"log/slog"

	"github.com/rs/zerolog"
	"github.com/topfreegames/pitaya/v2/logger/interfaces"
)

// CoreLogger CoreLogger
type CoreLogger struct {
	logger *slog.Logger
}

// NewCoreLogger NewCoreLogger
func NewCoreLogger(logger *slog.Logger) *CoreLogger {
	return &CoreLogger{
		logger: logger,
	}
}

// GetLogger GetLogger
func (l *CoreLogger) GetLogger() *slog.Logger {
	return l.logger
}

// Fatal Fatal
func (l *CoreLogger) Fatal(format ...any) {
	msg := fmt.Sprint(format...)
	l.logger.Error(msg)
	panic(msg)
}

// Fatalf Fatalf
func (l *CoreLogger) Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Error(msg)
	panic(msg)
}

// Fatalln Fatalln
func (l *CoreLogger) Fatalln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Error(msg)
	panic(msg)
}

// Debug Debug
func (l *CoreLogger) Debug(args ...any) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Debugf Debugf
func (l *CoreLogger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Debugln Debugln
func (l *CoreLogger) Debugln(args ...any) {
	l.logger.Debug(fmt.Sprintln(args...))
}

// Error Error
func (l *CoreLogger) Error(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
}

// Errorf Errorf
func (l *CoreLogger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Errorln Errorln
func (l *CoreLogger) Errorln(args ...any) {
	l.logger.Error(fmt.Sprintln(args...))
}

// Info Info
func (l *CoreLogger) Info(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

// Infof Infof
func (l *CoreLogger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Infoln Infoln
func (l *CoreLogger) Infoln(args ...any) {
	l.logger.Info(fmt.Sprintln(args...))
}

// Warn Warn
func (l *CoreLogger) Warn(args ...any) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warnf Warnf
func (l *CoreLogger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Warnln Warnln
func (l *CoreLogger) Warnln(args ...any) {
	l.logger.Warn(fmt.Sprintln(args...))
}

// Panic Panic
func (l *CoreLogger) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	l.logger.Error(msg)
	panic(msg)
}

// Panicf Panicf
func (l *CoreLogger) Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Error(msg)
	panic(msg)
}

// Panicln Panicln
func (l *CoreLogger) Panicln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Error(msg)
	panic(msg)
}

// WithFields WithFields
func (l *CoreLogger) WithFields(fields map[string]any) interfaces.Logger {
	logger := l.GetLogger()
	for k, v := range fields {
		logger = logger.With(k, v)
	}

	return NewCoreLogger(logger)
}

// WithField WithField
func (l *CoreLogger) WithField(key string, value any) interfaces.Logger {
	logger := l.logger.With(key, value)

	return NewCoreLogger(logger)
}

// WithError WithError
func (l *CoreLogger) WithError(err error) interfaces.Logger {
	logger := l.logger.With(zerolog.ErrorFieldName, err)

	return NewCoreLogger(logger)
}

// GetInternalLogger GetInternalLogger
func (l *CoreLogger) GetInternalLogger() any {
	return l.logger
}
