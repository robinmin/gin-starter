package config

import (
	"fmt"

	"github.com/robinmin/gin-starter/pkg/utility"
)

const AppName = "dragondoat"
const AppVersion = "0.0.1"

var __cfgFile string // global config file name

func Setup(cfgFile string) {
	__cfgFile = cfgFile
}

type AppConfig struct {
	System struct {
		DebugMode          bool   `yaml:"debug_mode,omitempty"`
		ServerAddr         string `yaml:"server_address,omitempty"`
		EnableCORS         bool   `yaml:"enable_cors,omitempty"`
		EnableAuth         bool   `yaml:"enable_auth,omitempty"`
		ExternalSvrAddress string `yaml:"external_svr_address,omitempty"`
	} `yaml:"system"`

	Database struct {
		Type     string `yaml:"dbtype,omitempty"`
		Host     string `yaml:"dbhost,omitempty"`
		Port     int    `yaml:"dbport,omitempty"`
		Database string `yaml:"dbname,omitempty"`
		User     string `yaml:"dbuser,omitempty"`
		Password string `yaml:"dbpassword,omitempty"`
	} `yaml:"database"`
}

func NewAppConfig() *AppConfig {
	cfg, err := utility.LoadConfig[AppConfig](__cfgFile)
	if err != nil {
		fmt.Println("Failed to load yaml config file from " + __cfgFile + "[" + err.Error() + "]")
		return &AppConfig{}
	}

	return cfg
}

func (cfg *AppConfig) SaveConfig() error {
	return utility.SaveConfig(cfg, __cfgFile)
}

func (cfg *AppConfig) GetConnectionString() (string, error) {
	switch cfg.Database.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database), nil
	case "sqlite3":
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", cfg.Database.Database), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
}
