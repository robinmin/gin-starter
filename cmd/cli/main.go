package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/robinmin/gin-starter/config"
	"github.com/robinmin/gin-starter/pkg/bootstrap"
	sloggin "github.com/samber/slog-gin"
	"go.uber.org/fx"
)

var (
	help        bool
	config_file string
	verbose     bool
)

func init() {
	// setup flags
	flag.BoolVar(&help, "h", false, "show the help message")
	flag.StringVar(&config_file, "f", "", "config file")
	flag.BoolVar(&verbose, "v", false, "show detail information")
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func NewAppConfig() *config.AppConfig {
	if config_file == "" {
		config_file = "config/app_config.yaml"
	}
	cfg, err := bootstrap.LoadConfig[config.AppConfig](config_file)
	if err != nil {
		fmt.Println("Failed to load yaml config file: " + err.Error())
		return nil
	}
	return cfg
}

func main() {
	// parse command line arguments and show help only if specified
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	app := fx.New(
		// configurations for logger and config file items
		fx.Provide(NewAppConfig),
		fx.Provide(func(cfg *config.AppConfig) *bootstrap.ApplicationConfig {
			return &bootstrap.ApplicationConfig{
				TrustedProxies: cfg.System.TrustedProxies,
				ServerAddr:     cfg.System.ServerAddr,
				Verbose:        verbose,
			}
		}),
		fx.Provide(func(cfg *config.AppConfig) bootstrap.LoggerParams {
			return bootstrap.LoggerParams{
				LogFileName:  fmt.Sprintf("log/gin-starter-%s.log", time.Now().Format("20060102")),
				DefaultLevel: slog.LevelDebug,
				Config: sloggin.Config{
					WithSpanID:  true,
					WithTraceID: true,
				},
			}
		}),

		// enable inported modules
		bootstrap.Module,

		// run application
		fx.Invoke(func(app *bootstrap.Application, svr *http.Server, logger *slog.Logger) {
			logger.Info("Application started")
		}),
	)
	app.Run()
}
