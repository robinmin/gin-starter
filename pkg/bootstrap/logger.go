package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.uber.org/fx"
)

var (
	__logFileHandler *os.File
)

type LoggerParams struct {
	// fx.In

	LogFileName  string
	DefaultLevel slog.Level
	Config       sloggin.Config
}

type AppLogger struct {
	// fx.In
	Log *slog.Logger

	Params LoggerParams
}

func (logger *AppLogger) Debug(msg string, args ...any) {
	logger.Log.Debug(msg, args...)
}

func (logger *AppLogger) Info(msg string, args ...any) {
	logger.Log.Info(msg, args...)
}

func (logger *AppLogger) Error(msg string, args ...any) {
	logger.Log.Error(msg, args...)
}

func (logger *AppLogger) Warn(msg string, args ...any) {
	logger.Log.Warn(msg, args...)
}

func (logger *AppLogger) Handler() slog.Handler {
	return logger.Log.Handler()
}

// func (logger *AppLogger) Enabled(ctx context.Context, level slog.Level) bool
// func (logger *AppLogger) Log(ctx context.Context, level slog.Level, msg string, args ...any)
// func (logger *AppLogger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
// func (logger *AppLogger) With(args ...any) *slog.Logger
// func (logger *AppLogger) WithGroup(name string) *slog.Logger
// func (logger *AppLogger) DebugContext(ctx context.Context, msg string, args ...any)
// func (logger *AppLogger) InfoContext(ctx context.Context, msg string, args ...any)
// func (logger *AppLogger) WarnContext(ctx context.Context, msg string, args ...any)
// func (logger *AppLogger) ErrorContext(ctx context.Context, msg string, args ...any)

func createLogWriter(filename string) io.Writer {
	var writers []io.Writer
	if gin.IsDebugging() {
		writers = append(writers, os.Stdout)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current working directory: ", cwd)

	if __logFileHandler == nil {
		var err error
		__logFileHandler, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
		if err != nil {
			fmt.Println("Failed to open log file: %w", err)
			return nil
		}
	}

	writers = append(writers, __logFileHandler)
	return io.MultiWriter(writers...)
}

func closeLogFile() {
	if __logFileHandler != nil {
		__logFileHandler.Close()
		__logFileHandler = nil
	}
}

func NewLogger(params LoggerParams, lc fx.Lifecycle) *AppLogger {
	opts := &slog.HandlerOptions{
		Level: params.DefaultLevel,
	}
	writer := createLogWriter(params.LogFileName)
	if writer == nil {
		gin.DefaultWriter = os.Stdout
	} else {
		gin.DefaultWriter = writer
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			closeLogFile()
			return nil
		},
	})
	return &AppLogger{
		Log:    slog.New(slog.NewTextHandler(gin.DefaultWriter, opts)),
		Params: params,
	}
}
