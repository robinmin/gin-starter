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

var __logFileHandler *os.File

type LoggerParams struct {
	// fx.In

	LogFileName  string
	DefaultLevel slog.Level
	Config       sloggin.Config
}

type AppLogger struct {
	// fx.In
	*slog.Logger

	Params LoggerParams
}

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
		Logger: slog.New(slog.NewTextHandler(gin.DefaultWriter, opts)),
		Params: params,
	}
}
