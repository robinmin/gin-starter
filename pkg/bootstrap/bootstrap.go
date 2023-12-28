package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	status "github.com/appleboy/gin-status-api"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
var (
	__errorInfo map[int]string
)

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

	// server instance
	server *http.Server

	// DB instance
	// DB     *database.DB

	lifeCycle fx.Lifecycle
}

func NewApplication(logger *AppLogger, cfg *ApplicationConfig, lc fx.Lifecycle) *Application {
	app := &Application{
		Config: cfg,
	}

	app.engine = gin.New()
	// The middleware will log all requests attributes.
	app.engine.Use(sloggin.NewWithConfig(logger.Logger, logger.Params.Config), gin.Recovery())
	app.engine.ForwardedByClientIP = true
	app.engine.Use(GlobalErrorMiddleware())

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

	app.server = NewHttpServer(app, logger)
	app.lifeCycle = lc

	// default status api
	app.engine.GET("/status", status.GinHandler)

	return app
}

func NewHttpServer(app *Application, logger *AppLogger) *http.Server {
	return &http.Server{
		Addr:         app.Config.ServerAddr,
		Handler:      app.engine,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
}

func SetErrorInfo(ei map[int]string) {
	__errorInfo = ei
}

func GetMessage(code int) string {
	if msg, ok := __errorInfo[code]; ok {
		return msg
	} else {
		return fmt.Sprintf("Unknown error code: %d", code)
	}
}

func (app *Application) RunServer(logger *AppLogger) error {
	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- app.server.Shutdown(ctx)
	}()

	app.lifeCycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting server", slog.Group("server", "addr", app.server.Addr))

			go func() {
				err := app.server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					shutdownErrorChan <- err
				}
			}()

			logger.Info("Succeeded to start HTTP Server at", slog.Group("server", "addr", app.server.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopped server", slog.Group("server", "addr", app.server.Addr))

			err := <-shutdownErrorChan
			if err != nil {
				return err
			}

			return app.server.Shutdown(ctx)
		},
	})
	return nil
}

func NewQuickResult(code int, data interface{}) Result {
	return Result{
		Code:    code,
		Message: GetMessage(code),
		Data:    data,
	}
}

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行请求
		c.Next()

		// 发生了错误
		if len(c.Errors) > 0 {
			// 获取最后一个error 返回
			err := c.Errors.Last()
			NewResult(http.StatusInternalServerError, err.Error(), nil).Fail(c)
			return
		}
	}
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

	err = os.WriteFile(yamlFile, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}
