package main

import (
	"flag"
	"fmt"

	"github.com/robinmin/gin-starter/config"
	"github.com/robinmin/gin-starter/pkg/bootstrap"
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
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
func newMyAppConfig() *config.MyAppConfig {
	if config_file == "" {
		config_file = "config/app_config.yaml"
	}
	cfg, err := bootstrap.LoadConfig[config.MyAppConfig](config_file)
	if err != nil {
		fmt.Println("Failed to load yaml config file: " + err.Error())
		return nil
	}

	return cfg
}

func main() {
	// ctx := context.Background()
	// Initialize logger.
	cleanLoggerFn, err := bootstrap.InitLogger()
	if err != nil {
		panic(err)
	}
	defer cleanLoggerFn()

	// parse command line arguments and show help only if specified
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	// set error information
	bootstrap.SetErrorInfo(config.ErrorCodeMapping)

	fx.New(
		fx.WithLogger(func(log *bootstrap.AppLogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Logger}
		}),
		// configurations for logger and config file items
		fx.Provide(newMyAppConfig),
		fx.Provide(func(cfg *config.MyAppConfig) types.AppConfig {
			sc := cfg.Basic
			sc.Sentry.EventsMeta = config.SentryEventsMeta
			// sc.Log.Config = sloggin.Config{
			// 	WithSpanID:  true,
			// 	WithTraceID: true,
			// }
			return sc
		}),

		// enable inported modules
		bootstrap.Module,

		// run application
		fx.Invoke(func(app *bootstrap.Application, logger *bootstrap.AppLogger) {
			if err := app.RunServer(logger); err != nil {
				logger.Error("Failed to run server : " + err.Error())
			} else {
				logger.Info("Succeeded to run server")
			}
		}),
	).Run()
}
