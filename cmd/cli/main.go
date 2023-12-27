package main

import (
	"flag"
	"fmt"
	"log/slog"
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

	svr := fx.New(
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
		fx.Provide(func(cfg *config.AppConfig) bootstrap.DBParams {
			return bootstrap.DBParams{
				Type:     cfg.Database.Type,
				Host:     cfg.Database.Host,
				Port:     cfg.Database.Port,
				Database: cfg.Database.Database,
				User:     cfg.Database.User,
				Password: cfg.Database.Password,
			}
		}),

		// enable inported modules
		bootstrap.Module,

		// run application
		fx.Invoke(func(app *bootstrap.Application, logger *slog.Logger) {
			if err := app.RunServer(logger); err != nil {
				logger.Error("Failed to run server : " + err.Error())
			} else {
				logger.Info("Succeeded to run server")
			}
		}),
	)
	svr.Run()
}
