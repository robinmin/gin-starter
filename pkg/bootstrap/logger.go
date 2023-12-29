package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
	"go.uber.org/fx"
)

type LoggerParams types.AppLogConfig

type AppLogger struct {
	*slog.Logger // Logger instance

	Params LoggerParams // Parameters for the logger
}

// var __logFileHandler *os.File

func (lp *LoggerParams) CreateLogWriter(filename string) io.Writer {
	var writers []io.Writer
	if gin.IsDebugging() {
		writers = append(writers, os.Stdout)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current working directory: ", cwd)

	if lp.LogFileHandler == nil {
		var err error
		lp.LogFileHandler, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
		if err != nil {
			fmt.Println("Failed to open log file: %w", err)
			return nil
		}
	}

	writers = append(writers, lp.LogFileHandler)
	return io.MultiWriter(writers...)
}

func (lp *LoggerParams) CloseLogFile() {
	if lp.LogFileHandler != nil {
		lp.LogFileHandler.Close()
		lp.LogFileHandler = nil
	}
}

func NewLogger(cfg types.AppConfig, lc fx.Lifecycle) *AppLogger {
	params := LoggerParams(cfg.Log)
	opts := &slog.HandlerOptions{
		Level: slog.Level(params.DefaultLevel),
	}
	writer := params.CreateLogWriter(params.LogFileName)
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
			params.CloseLogFile()
			return nil
		},
	})

	return &AppLogger{
		Logger: slog.New(slog.NewTextHandler(gin.DefaultWriter, opts)),
		Params: params,
	}
}
