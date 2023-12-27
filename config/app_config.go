package config

const (
	AppName    = "dragondoat"
	AppVersion = "0.0.1"
)

type AppConfig struct {
	System struct {
		DebugMode          bool   `yaml:"debug_mode,omitempty"`
		ServerAddr         string `yaml:"server_address,omitempty"`
		EnableCORS         bool   `yaml:"enable_cors,omitempty"`
		EnableAuth         bool   `yaml:"enable_auth,omitempty"`
		ExternalSvrAddress string `yaml:"external_svr_address,omitempty"`
		TrustedProxies     string `yaml:"trusted_proxies,omitempty"`
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
