package bootstrap

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"os"

	status "github.com/appleboy/gin-status-api"
	sloggin "github.com/samber/slog-gin"
	"go.uber.org/fx"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
type LoggerParams struct {
	// fx.In

	LogFileName  string
	DefaultLevel slog.Level
	Config       sloggin.Config
}

var __logFileHandler *os.File

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
		__logFileHandler, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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

func NewLogger(params LoggerParams, lc fx.Lifecycle) *slog.Logger {
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
	return slog.New(slog.NewTextHandler(gin.DefaultWriter, opts))
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type ApplicationConfig struct {
	// Configuration of server trusted proxies
	TrustedProxies string

	// Configuration of server address
	ServerAddr string

	// verbose flag
	Verbose bool
}

type Application struct {
	// Configuration
	Config *ApplicationConfig

	// Engine instance
	engine *gin.Engine

	// DB instance
	// DB     *database.DB
}

func NewApplication(logger *slog.Logger, logParam LoggerParams, cfg *ApplicationConfig) *Application {
	app := &Application{
		Config: cfg,
	}

	app.engine = gin.New()
	// The middleware will log all requests attributes.
	app.engine.Use(sloggin.NewWithConfig(logger, logParam.Config), gin.Recovery())
	app.engine.ForwardedByClientIP = true

	var err error
	if cfg.TrustedProxies == "" {
		err = app.engine.SetTrustedProxies([]string{"127.0.0.1"})
	} else {
		err = app.engine.SetTrustedProxies(strings.Split(cfg.TrustedProxies, ";"))
	}
	if err != nil {
		logger.Warn("Failed to set trusted proxies")
	}

	if gin.IsDebugging() {
		gin.ForceConsoleColor()
	} else {
		// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
		gin.DisableConsoleColor()
	}

	// default status api
	app.engine.GET("/status", status.GinHandler)

	return app
}

type HttpServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	App       *Application
}

func NewHttpServer(params HttpServerParams, logger *slog.Logger) *http.Server {
	srv := &http.Server{
		Addr:         params.App.Config.ServerAddr,
		Handler:      params.App.engine,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting server", slog.Group("server", "addr", srv.Addr))

			go func() {
				err := srv.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					shutdownErrorChan <- err
				}
			}()

			logger.Info("Succeeded to start HTTP Server at", slog.Group("server", "addr", srv.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopped server", slog.Group("server", "addr", srv.Addr))

			err := <-shutdownErrorChan
			if err != nil {
				return err
			}

			return srv.Shutdown(ctx)
		},
	})

	return srv
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

// LoadConfig 从指定的YAML文件中加载配置信息
func LoadConfig[T any](yamlFile string) (*T, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var config T
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 将配置信息保存到指定的YAML文件中
func SaveConfig[T any](cfg *T, yamlFile string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(yamlFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// func (cfg *AppConfig) GetConnectionString() (string, error) {
// 	switch cfg.Database.Type {
// 	case "mysql":
// 		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database), nil
// 	case "sqlite3":
// 		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", cfg.Database.Database), nil
// 	default:
// 		return "", fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
// 	}
// }
