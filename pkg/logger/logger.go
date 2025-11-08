package logger

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"silo/pkg/logger/config"
	"silo/pkg/logger/handler"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/topfreegames/pitaya/v2/constants"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log Log
var Log *slog.Logger
var LogWriter io.Writer
var LogLevel slog.Level

// InitDefaultLogger default logger
func InitDefaultLogger(cfg *config.LogConfig) {
	initLevel(cfg)
	initLogWriter(cfg)
	initLoggerWithConfig(cfg)
}

func initLevel(cfg *config.LogConfig) {
	slv := slog.LevelInfo
	if err := (&slv).UnmarshalText([]byte(cfg.Level)); err != nil {
		panic(fmt.Sprintf("invalid log level: %s", cfg.Level))
	}

	// set global
	LogLevel = slv
}

func initLogWriter(cfg *config.LogConfig) {
	// writers
	var writers []io.Writer
	if cfg.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05.000",
			NoColor:    cfg.NoColor,
		})
	}
	if cfg.FileLoggingEnabled {
		writers = append(writers, newRollingFile(cfg))
	}
	mw := io.MultiWriter(writers...)

	// set global
	LogWriter = mw
}

// NewLoggerWithConfig NewLoggerWithConfig
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func initLoggerWithConfig(cfg *config.LogConfig) {
	// level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		panic(fmt.Sprintf("invalid log level: %s", cfg.Level))
	}

	// global level
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = cfg.TimeFieldFormat
	zerolog.TimestampFieldName = cfg.TimestampFieldName
	zerolog.LevelFieldName = cfg.LevelFieldName
	zerolog.MessageFieldName = cfg.MessageFieldName
	zerolog.ErrorFieldName = cfg.ErrorFieldName
	// use utc time
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	// log level to upper case
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}

	zlogger := zerolog.New(LogWriter)
	logger := slog.New(handler.Option{
		Level:     LogLevel,
		Logger:    &zlogger,
		StructLog: cfg.StructLog,
		Pretty:    cfg.Pretty,
	}.NewZerologHandler())

	if cfg.WithServeInfo {
		// logger.Handler().WithGroup()
		logger = logger.With(slog.Group("server", "type", cfg.ServerType, "id", cfg.ServerID))
	}

	// set global
	Log = logger
	slog.SetDefault(logger)
}

// GetLoggerFromCtx returns the default logger from the given context
func GetLoggerFromCtx(ctx context.Context) *slog.Logger {
	l := ctx.Value(constants.LoggerCtxKey)
	if l == nil {
		return Log
	}
	coreLogger := l.(*CoreLogger)

	return coreLogger.GetLogger()
}

func newRollingFile(cfg *config.LogConfig) io.Writer {
	if err := os.MkdirAll(cfg.Directory, 0744); err != nil {
		panic(fmt.Sprintf("can't create log directory %v, error %v", cfg.Directory, err))
	}

	fileName := cfg.Filename
	if fileName == "" {
		fileName = fmt.Sprintf("%s.log", cfg.ServerType)
	}

	// create default log file with permission
	// creating the original file yourself with the permissions you want, and then lumberjack will respect those and copy them to any new file it creates
	// https://github.com/natefinch/lumberjack/issues/82
	filePath := path.Join(cfg.Directory, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 8 is the traditional base for unix file modes
		mode, err := strconv.ParseUint(cfg.FilePermission, 8, 32)
		if err != nil {
			panic(fmt.Sprintf("fail to parse FilePermission: %s", err.Error()))
		}
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, fs.FileMode(mode))
		if err != nil {
			panic(fmt.Sprintf("fail to create log file: %s", err.Error()))
		}
		defer file.Close()
	}

	return &lumberjack.Logger{
		Filename:   filePath,
		MaxBackups: cfg.MaxBackups, // files
		MaxSize:    cfg.MaxSize,    // megabytes
		MaxAge:     cfg.MaxAge,     // days
	}
}
