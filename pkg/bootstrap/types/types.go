package types

import (
	"log/slog"
	"os"
	"time"

	sloggin "github.com/samber/slog-gin"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Definitions for sentry
type UserDefinedEvent int
type UserDefinedEventMeta struct {
	Name  string     // Name of the event
	Level slog.Level // Level of the event level
	Group string     // Group of the event belongs to
}
type UserDefinedEventMap map[UserDefinedEvent]UserDefinedEventMeta

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Definitions for system configuration
type AppSysConfig struct {
	DebugMode          bool   `yaml:"debug_mode,omitempty" json:"debug_mode,omitempty" default:"false"`
	ServerAddr         string `yaml:"server_address,omitempty" json:"server_address,omitempty" default:":7086"`
	EnableCORS         bool   `yaml:"enable_cors,omitempty" json:"enable_cors,omitempty" default:"true"`
	EnableAuth         bool   `yaml:"enable_auth,omitempty" json:"enable_auth,omitempty" default:"true"`
	ExternalSvrAddress string `yaml:"external_svr_address,omitempty" json:"external_svr_address,omitempty" default:""`
	TrustedProxies     string `yaml:"trusted_proxies,omitempty" json:"trusted_proxies,omitempty" default:"127.0.0.1;10.0.0.0/8"`
	StaticDir          string `yaml:"static_dir,omitempty" json:"static_dir,omitempty" default:"./static"`
	StaticURL          string `yaml:"static_url,omitempty" json:"static_url,omitempty" default:"/static"`
}

type AppLogConfig struct {
	LogFileName       string         `yaml:"log_file_name,omitempty" json:"log_file_name,omitempty" default:"app.log"`                // Log file name
	LogFileNameFormat string         `yaml:"log_file_name_format,omitempty" json:"log_file_name_format,omitempty" default:"20060102"` // Log file name format
	DefaultLevel      int            `yaml:"default_level,omitempty" json:"default_level,omitempty" default:"-4"`                     // Default level of the logger
	Config            sloggin.Config `yaml:"-" json:"-"`                                                                              // Configuration of the logger
	LogFileHandler    *os.File       `yaml:"-" json:"-"`                                                                              // Handler for the log file
}

// Definitions for database configuration
type AppDBConfig struct {
	Type     string `yaml:"dbtype,omitempty" json:"dbtype,omitempty" default:"sqlite3"`
	Host     string `yaml:"dbhost,omitempty" json:"dbhost,omitempty" default:"localhost"`
	Port     int    `yaml:"dbport,omitempty" json:"dbport,omitempty" default:"3306"`
	Database string `yaml:"dbname,omitempty" json:"dbname,omitempty" default:"database"`
	User     string `yaml:"dbuser,omitempty" json:"dbuser,omitempty" default:"user"`
	Password string `yaml:"dbpassword,omitempty" json:"dbpassword,omitempty" default:""`
}

// Definitions for redis configuration
type AppRedisConfig struct {
	Size               int           `yaml:"size,omitempty" json:"size,omitempty" default:"10"`                                   // maximum number of idle connections.
	Network            string        `yaml:"network,omitempty" json:"network,omitempty" default:"tcp"`                            // tcp or udp
	Address            string        `yaml:"address,omitempty" json:"address,omitempty" default:"localhost:6379"`                 // host:port of redis server
	Password           string        `yaml:"password,omitempty" json:"password,omitempty" default:""`                             // redis-password
	DB                 string        `yaml:"db,omitempty" json:"db,omitempty" default:"0"`                                        // database
	KeyPairs           string        `yaml:"key_pairs,omitempty" json:"key_pairs,omitempty" default:""`                           // Keys are defined in pairs to allow key rotation, but the common case is to set a single authentication key and optionally an encryption key.
	DefaultExpiration  time.Duration `yaml:"default_expiration,omitempty" json:"default_expiration,omitempty" default:"10m"`      // default expiration time for redis cache
	EnableRedisCache   bool          `yaml:"enable_redis_cache,omitempty" json:"enable_redis_cache,omitempty" default:"true"`     // use redis cache
	EnableRedisSession bool          `yaml:"enable_redis_session,omitempty" json:"enable_redis_session,omitempty" default:"true"` // use redis session
	SessionName        string        `yaml:"session_name,omitempty" json:"session_name,omitempty" default:"session"`              // session name
}

// Definitions for sentry configuration
type AppSentryConfig struct {
	DSN              string              `yaml:"sentry_dsn,omitempty" json:"sentry_dsn,omitempty" default:""`                    // DSN of the sentry
	TracesSampleRate float64             `yaml:"traces_sample_rate,omitempty" json:"traces_sample_rate,omitempty" default:"1.0"` // trac sample rate
	DefaultLevel     int                 `yaml:"default_level,omitempty" json:"default_level,omitempty" default:"-4"`            // Default level of the sentry
	EventsMeta       UserDefinedEventMap `yaml:"-" json:"-"`                                                                     // Events meatadata mappings
}

type AppConfig struct {
	System   AppSysConfig    `yaml:"system,omitempty" json:"system,omitempty"`
	Log      AppLogConfig    `yaml:"log,omitempty" json:"log,omitempty"`
	Database AppDBConfig     `yaml:"database,omitempty" json:"database,omitempty"`
	Redis    AppRedisConfig  `yaml:"redis,omitempty" json:"redis,omitempty"`
	Sentry   AppSentryConfig `yaml:"sentry,omitempty" json:"sentry,omitempty"`
}
