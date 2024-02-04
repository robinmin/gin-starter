package bootstrap

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	status "github.com/appleboy/gin-status-api"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/creasty/defaults"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
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

func NewApplication(
	cfg types.AppConfig,
	lc fx.Lifecycle,
	sty *AppSentry,
	logger *AppLogger,
) *Application {
	app := &Application{
		Config: cfg.System,
	}

	app.engine = gin.New()
	if gin.IsDebugging() {
		gin.ForceConsoleColor()
	} else {
		// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
		gin.DisableConsoleColor()
	}
	app.engine.ForwardedByClientIP = true

	var err error
	if cfg.System.TrustedProxies == "" {
		err = app.engine.SetTrustedProxies([]string{"127.0.0.1"})
	} else {
		err = app.engine.SetTrustedProxies(strings.Split(cfg.System.TrustedProxies, ";"))
	}
	if err != nil {
		logger.Warn("Failed to set trusted proxies")
	}

	// The middleware functions are executed in the order they are defined.
	if err = app.useMiddlewares(context.Background(), cfg, logger); err != nil {
		logger.Error("Failed to enable all middlewares: " + err.Error())
	}

	app.server = NewHttpServer(app, logger)
	app.lifeCycle = lc

	// default status api
	app.engine.GET("/status", status.GinHandler)

	return app
}

func (app *Application) useMiddlewares(ctx context.Context, cfg types.AppConfig, logger *AppLogger) error {
	// global middlewares for error handling
	app.engine.Use(GlobalErrorHandler())

	// Middleware for logging
	app.engine.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: cfg.Middlewares.Log.TimeFormat,
		UTC:        cfg.Middlewares.Log.UTC,
		SkipPaths:  cfg.Middlewares.Log.SkipPaths,
	}))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	app.engine.Use(ginzap.RecoveryWithZap(logger, true))

	// Middleware for CORS
	if cfg.Middlewares.CORS.Enable {
		app.engine.Use(cors.New(cors.Config{
			AllowOriginFunc: func(origin string) bool { return true },

			AllowAllOrigins:           cfg.Middlewares.CORS.AllowAllOrigins,
			AllowMethods:              cfg.Middlewares.CORS.AllowMethods,
			AllowPrivateNetwork:       cfg.Middlewares.CORS.AllowPrivateNetwork,
			AllowHeaders:              cfg.Middlewares.CORS.AllowHeaders,
			AllowCredentials:          cfg.Middlewares.CORS.AllowCredentials,
			ExposeHeaders:             cfg.Middlewares.CORS.ExposeHeaders,
			MaxAge:                    time.Second * time.Duration(cfg.Middlewares.CORS.MaxAge),
			AllowWildcard:             cfg.Middlewares.CORS.AllowWildcard,
			AllowBrowserExtensions:    cfg.Middlewares.CORS.AllowBrowserExtensions,
			AllowWebSockets:           cfg.Middlewares.CORS.AllowWebSockets,
			AllowFiles:                cfg.Middlewares.CORS.AllowFiles,
			OptionsResponseStatusCode: cfg.Middlewares.CORS.OptionsResponseStatusCode,
		}))
	}

	// Middleware for session
	if cfg.Middlewares.Session.Enable {
		if cfg.Middlewares.Session.UseRedis {
			rstore, _ := redis.NewStoreWithDB(
				cfg.Redis.Size,
				cfg.Redis.Network,
				cfg.Redis.Address,
				cfg.Redis.Password,
				cfg.Redis.DB,
				[]byte(cfg.Redis.KeyPairs))
			app.engine.Use(sessions.Sessions(cfg.Middlewares.Session.Name, rstore))
		} else {
			app.engine.Use(sessions.Sessions(cfg.Middlewares.Session.Name, cookie.NewStore([]byte(cfg.Redis.KeyPairs))))
		}
	}

	// Middleware for static files
	if cfg.Middlewares.Static.Enable && len(cfg.Middlewares.Static.StaticDir) > 0 && len(cfg.Middlewares.Static.StaticURL) > 0 {
		if absPath, err := filepath.Abs(cfg.Middlewares.Static.StaticDir); err == nil {
			logger.Debug("URL " + cfg.Middlewares.Static.StaticURL + " -> " + absPath)
			app.engine.Use(static.Serve(cfg.Middlewares.Static.StaticURL, static.LocalFile(cfg.Middlewares.Static.StaticDir, cfg.Middlewares.Static.Indexes)))
		} else {
			logger.Error("Failed to get absolute path of " + cfg.Middlewares.Static.StaticDir)
			return err
		}
	}

	if cfg.Middlewares.Gzip.Enable {
		app.engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	return nil
}

func NewRedisCache(cfg types.AppConfig) *persistence.RedisStore {
	if cfg.Middlewares.Cache.Enable && cfg.Middlewares.Cache.UseRedis {
		return persistence.NewRedisCache(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DefaultExpiration)
	}
	return nil
}

func NewMemoryCache(cfg types.AppConfig) *persistence.InMemoryStore {
	if cfg.Middlewares.Cache.Enable && !cfg.Middlewares.Cache.UseRedis {
		return persistence.NewInMemoryStore(cfg.Redis.DefaultExpiration)
	}
	return nil
}

func NewHttpServer(app *Application, logger *AppLogger) *http.Server {
	return &http.Server{
		Addr:         app.Config.ServerAddr,
		Handler:      app.engine,
		ErrorLog:     logger.GetRawLogger(),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
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
			logger.Info("Starting server", zap.String("server address", app.server.Addr))

			go func() {
				err := app.server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					shutdownErrorChan <- err
				}
			}()

			logger.Info("Succeeded to start HTTP Server at", zap.String("server address", app.server.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopped server", zap.String("server address", app.server.Addr))

			err := <-shutdownErrorChan
			if err != nil {
				return err
			}

			return app.server.Shutdown(ctx)
		},
	})
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// For simulating ternary expressions that golang lacks
// func ifelse[T any](condition bool, true_part T, false_part T) T {
// 	if condition {
// 		return true_part
// 	}
// 	return false_part
// }

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
