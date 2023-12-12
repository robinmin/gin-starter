package main

import (
	"fmt"
	"os"

	"github.com/robinmin/gin-starter/config"
	"github.com/robinmin/gin-starter/pkg/bootstrap"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	cfg, err0 := bootstrap.LoadConfig[config.AppConfig]("config/app_config.yaml")
	if err0 != nil {
		fmt.Println("Failed to load yaml config file: " + err0.Error())
		os.Exit(1)
	}

	app, err := bootstrap.NewApplication[config.AppConfig](cfg, "log/gin-starter.log")
	if err != nil {
		fmt.Println("Failed to create an application instance on startup: " + err.Error())
		os.Exit(1)
	}
	defer app.Quit()

	err = app.RunServer(app.Config.System.ServerAddr)
	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
}
