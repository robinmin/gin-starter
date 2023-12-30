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
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"

	"github.com/creasty/defaults"
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
type Application struct {
	// Configuration
	Config types.AppSysConfig

	// Engine instance
	engine *gin.Engine

	// server instance
	server *http.Server

	// DB instance
	// DB     *database.DB

	lifeCycle fx.Lifecycle
}

func NewApplication(cfg types.AppConfig, lc fx.Lifecycle, sty *AppSentry, logger *AppLogger) *Application {
	app := &Application{
		Config: cfg.System,
	}

	app.engine = gin.New()
	// The middleware will log all requests attributes.
	app.engine.Use(sloggin.NewWithConfig(logger.Logger, logger.Params.Config), gin.Recovery())
	app.engine.ForwardedByClientIP = true
	app.engine.Use(GlobalErrorMiddleware())

	if cfg.System.EnableCORS {
		app.engine.Use(cors.New(cors.Config{
			AllowCredentials: true,
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowHeaders:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		}))
	}

	var err error
	if cfg.System.TrustedProxies == "" {
		err = app.engine.SetTrustedProxies([]string{"127.0.0.1"})
	} else {
		err = app.engine.SetTrustedProxies(strings.Split(cfg.System.TrustedProxies, ";"))
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

	// static files
	if len(cfg.System.StaticDir) > 0 && len(cfg.System.StaticURL) > 0 {
		app.engine.Use(static.Serve(cfg.System.StaticURL, static.LocalFile(cfg.System.StaticDir, true)))

	}

	app.engine.Use(static.Serve("/", static.LocalFile("/tmp", false)))

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

// create instance and load default values which defined in the struct definition
func NewInstance[T any]() *T {
	obj := new(T)
	if err := defaults.Set(obj); err != nil {
		return nil
	}
	return obj
}

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
