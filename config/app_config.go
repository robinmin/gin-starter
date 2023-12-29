package config

import "github.com/robinmin/gin-starter/pkg/bootstrap/types"

const (
	AppName    = "gin-stater"
	AppVersion = "0.0.1"
)

type MyAppConfig struct {
	Basic types.AppConfig `yaml:"basic"`
}
