package logger

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/labstack/gommon/log"
)

// EchoLogger EchoLogger
type EchoLogger struct {
	logger *slog.Logger
	prefix string
}

// NewEchoLogger NewEchoLogger
func NewEchoLogger(logger *slog.Logger) *EchoLogger {
	return &EchoLogger{
		logger: logger,
	}
}

// Output Output
func (l *EchoLogger) Output() io.Writer {
	return LogWriter
}

// SetOutput SetOutput
func (l *EchoLogger) SetOutput(w io.Writer) {
}

// Prefix Prefix
func (l *EchoLogger) Prefix() string {
	return l.prefix
}

// SetPrefix SetPrefix
func (l *EchoLogger) SetPrefix(p string) {
	l.prefix = p
}

// Level Level
func (l *EchoLogger) Level() log.Lvl {
	level := LogLevel
	switch level {
	case slog.LevelDebug:
		return log.DEBUG
	case slog.LevelInfo:
		return log.INFO
	case slog.LevelWarn:
		return log.WARN
	case slog.LevelError:
		return log.ERROR
	}
	return log.DEBUG
}

// SetLevel SetLevel
func (l *EchoLogger) SetLevel(v log.Lvl) {
}

// SetHeader SetHeader
func (l *EchoLogger) SetHeader(h string) {
	l.logger = l.logger.With("header", h)
}

// Print Print
func (l *EchoLogger) Print(i ...any) {
	l.logger.Debug(fmt.Sprint(i...))
}

// Printf Printf
func (l *EchoLogger) Printf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Printj Printj
func (l *EchoLogger) Printj(j log.JSON) {
	l.logger.Debug("", "json", j)
}

// Debug Debug
func (l *EchoLogger) Debug(i ...any) {
	l.logger.Debug(fmt.Sprint(i...))
}

// Debugf Debugf
func (l *EchoLogger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Debugj Debugj
func (l *EchoLogger) Debugj(j log.JSON) {
	l.logger.Debug("", "json", j)
}

// Info Info
func (l *EchoLogger) Info(i ...any) {
	l.logger.Info(fmt.Sprint(i...))
}

// Infof Infof
func (l *EchoLogger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Infoj Infoj
func (l *EchoLogger) Infoj(j log.JSON) {
	l.logger.Info("", "json", j)
}

// Warn Warn
func (l *EchoLogger) Warn(i ...any) {
	l.logger.Warn(fmt.Sprint(i...))
}

// Warnf Warnf
func (l *EchoLogger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Warnj Warnj
func (l *EchoLogger) Warnj(j log.JSON) {
	l.logger.Warn("", "json", j)
}

// Error Error
func (l *EchoLogger) Error(i ...any) {
	l.logger.Error(fmt.Sprint(i...))
}

// Errorf Errorf
func (l *EchoLogger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Errorj Errorj
func (l *EchoLogger) Errorj(j log.JSON) {
	l.logger.Error("", "json", j)
}

// Fatal Fatal
func (l *EchoLogger) Fatal(i ...any) {
	l.logger.Error(fmt.Sprint(i...))
}

// Fatalj Fatalj
func (l *EchoLogger) Fatalj(j log.JSON) {
	l.logger.Error("", "json", j)
}

// Fatalf Fatalf
func (l *EchoLogger) Fatalf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Panic Panic
func (l *EchoLogger) Panic(i ...any) {
	l.logger.Error(fmt.Sprint(i...))
}

// Panicj Panicj
func (l *EchoLogger) Panicj(j log.JSON) {
	l.logger.Error("", "json", j)
}

// Panicf Panicf
func (l *EchoLogger) Panicf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}
