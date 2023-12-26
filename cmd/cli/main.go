package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/robinmin/gin-starter/config"
	"github.com/robinmin/gin-starter/pkg/bootstrap"
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
func main() {
	// parse command line arguments and show help only if specified
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	if config_file == "" {
		config_file = "config/app_config.yaml"
	}
	cfg, err0 := bootstrap.LoadConfig[config.AppConfig](config_file)
	if err0 != nil {
		fmt.Println("Failed to load yaml config file: " + err0.Error())
		os.Exit(1)
	}

	appCfg := &bootstrap.ApplicationConfig{
		LogFileName: fmt.Sprintf("log/gin-starter-%s.log", time.Now().Format("20060102")),
		Verbose:     verbose,
	}
	app, err := bootstrap.NewApplication(appCfg)
	if err != nil {
		fmt.Println("Failed to create an application instance on startup: " + err.Error())
		os.Exit(1)
	}
	defer app.Quit()

	err = app.RunServer(cfg.System.ServerAddr)
	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
}
