package types

import (
	"time"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Definitions for sentry
type UserDefinedEvent int
type UserDefinedEventMeta struct {
	Name  string // Name of the event
	Level string // Level of the event level
	Group string // Group of the event belongs to
}
type UserDefinedEventMap map[UserDefinedEvent]UserDefinedEventMeta

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Definitions for system configuration
type AppSysConfig struct {
	DebugMode          bool   `yaml:"debug_mode,omitempty" json:"debug_mode,omitempty" default:"false"`
	ServerAddr         string `yaml:"server_address,omitempty" json:"server_address,omitempty" default:":7086"`
	EnableAuth         bool   `yaml:"enable_auth,omitempty" json:"enable_auth,omitempty" default:"true"`
	ExternalSvrAddress string `yaml:"external_svr_address,omitempty" json:"external_svr_address,omitempty" default:""`
	TrustedProxies     string `yaml:"trusted_proxies,omitempty" json:"trusted_proxies,omitempty" default:"127.0.0.1;10.0.0.0/8"`
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
	Size              int           `yaml:"size,omitempty" json:"size,omitempty" default:"10"`                              // maximum number of idle connections.
	Network           string        `yaml:"network,omitempty" json:"network,omitempty" default:"tcp"`                       // tcp or udp
	Address           string        `yaml:"address,omitempty" json:"address,omitempty" default:"localhost:6379"`            // host:port of redis server
	Password          string        `yaml:"password,omitempty" json:"password,omitempty" default:""`                        // redis-password
	DB                string        `yaml:"db,omitempty" json:"db,omitempty" default:"0"`                                   // database
	KeyPairs          string        `yaml:"key_pairs,omitempty" json:"key_pairs,omitempty" default:""`                      // Keys are defined in pairs to allow key rotation, but the common case is to set a single authentication key and optionally an encryption key.
	DefaultExpiration time.Duration `yaml:"default_expiration,omitempty" json:"default_expiration,omitempty" default:"10m"` // default expiration time for redis cache
	// EnableRedisCache  bool          `yaml:"enable_redis_cache,omitempty" json:"enable_redis_cache,omitempty" default:"true"` // use redis cache
}

// Definitions for sentry configuration
type AppSentryConfig struct {
	DSN              string              `yaml:"sentry_dsn,omitempty" json:"sentry_dsn,omitempty" default:""`                    // DSN of the sentry
	TracesSampleRate float64             `yaml:"traces_sample_rate,omitempty" json:"traces_sample_rate,omitempty" default:"1.0"` // trac sample rate
	DefaultLevel     int                 `yaml:"default_level,omitempty" json:"default_level,omitempty" default:"-4"`            // Default level of the sentry
	EventsMeta       UserDefinedEventMap `yaml:"-" json:"-"`                                                                     // Events meatadata mappings
}

type AppConfig struct {
	System      AppSysConfig    `yaml:"system,omitempty" json:"system,omitempty"`
	Database    AppDBConfig     `yaml:"database,omitempty" json:"database,omitempty"`
	Redis       AppRedisConfig  `yaml:"redis,omitempty" json:"redis,omitempty"`
	Sentry      AppSentryConfig `yaml:"sentry,omitempty" json:"sentry,omitempty"`
	Middlewares struct {
		Log struct {
			TimeFormat   string   `yaml:"time_format,omitempty" json:"time_format,omitempty" default:"2006-01-02T15:04:05Z07:00"`
			UTC          bool     `yaml:"utc,omitempty" json:"utc,omitempty" default:"false"`
			SkipPaths    []string `yaml:"skip_paths,omitempty" json:"skip_paths,omitempty"`
			DefaultLevel string   `yaml:"default_level,omitempty" json:"default_level,omitempty" default:"info"` // Default level of the logger
		} `yaml:"log,omitempty" json:"log,omitempty"`

		CORS struct {
			Enable bool `yaml:"enable,omitempty" json:"enable,omitempty" default:"true"`

			AllowAllOrigins bool `yaml:"allow_all_origins,omitempty" json:"allow_all_origins,omitempty" default:"true"`

			// AllowOrigins is a list of origins a cross-domain request can be executed from.
			// If the special "*" value is present in the list, all origins will be allowed.
			// Default value is []
			AllowOrigins []string `yaml:"allow_origins,omitempty" json:"allow_origins,omitempty" default:""`

			// AllowOriginFunc is a custom function to validate the origin. It takes the origin
			// as an argument and returns true if allowed or false otherwise. If this option is
			// set, the content of AllowOrigins is ignored.
			// AllowOriginFunc func(origin string) bool

			// AllowMethods is a list of methods the client is allowed to use with
			// cross-domain requests. Default value is simple methods (GET, POST, PUT, PATCH, DELETE, HEAD, and OPTIONS)
			AllowMethods []string `yaml:"allow_methods,omitempty" json:"allow_methods,omitempty" default:"GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"`

			// AllowPrivateNetwork indicates whether the response should include allow private network header
			AllowPrivateNetwork bool `yaml:"allow_private_network,omitempty" json:"allow_private_network,omitempty" default:"false"`

			// AllowHeaders is list of non simple headers the client is allowed to use with
			// cross-domain requests.
			AllowHeaders []string `yaml:"allow_headers,omitempty" json:"allow_headers,omitempty" default:"Origin,Content-Length,Content-Type"`

			// AllowCredentials indicates whether the request can include user credentials like
			// cookies, HTTP authentication or client side SSL certificates.
			AllowCredentials bool `yaml:"allow_credentials,omitempty" json:"allow_credentials,omitempty" default:"false"`

			// ExposeHeaders indicates which headers are safe to expose to the API of a CORS
			// API specification
			ExposeHeaders []string `yaml:"expose_headers,omitempty" json:"expose_headers,omitempty" default:""`

			// MaxAge indicates how long (with second-precision) the results of a preflight request
			// can be cached
			MaxAge time.Duration `yaml:"max_age,omitempty" json:"max_age,omitempty"` // default as 12 * time.Hour

			// Allows to add origins like http://some-domain/*, https://api.* or http://some.*.subdomain.com
			AllowWildcard bool `yaml:"allow_wildcard,omitempty" json:"allow_wildcard,omitempty" default:"false"`

			// Allows usage of popular browser extensions schemas
			AllowBrowserExtensions bool `yaml:"allow_browser_extensions,omitempty" json:"allow_browser_extensions,omitempty" default:"false"`

			// Allows usage of WebSocket protocol
			AllowWebSockets bool `yaml:"allow_web_sockets,omitempty" json:"allow_web_sockets,omitempty" default:"false"`

			// Allows usage of file:// schema (dangerous!) use it only when you 100% sure it's needed
			AllowFiles bool `yaml:"allow_files,omitempty" json:"allow_files,omitempty" default:"false"`

			// Allows to pass custom OPTIONS response status code for old browsers / clients
			OptionsResponseStatusCode int `yaml:"options_response_status_code,omitempty" json:"options_response_status_code,omitempty" default:"200"`
		} `yaml:"cors,omitempty" json:"cors,omitempty"`

		Session struct {
			Enable   bool   `yaml:"enable,omitempty" json:"enable,omitempty" default:"true"`
			Name     string `yaml:"name,omitempty" json:"name,omitempty" default:"session"`        // session name
			UseRedis bool   `yaml:"use_redis,omitempty" json:"use_redis,omitempty" default:"true"` // use redis session
		} `yaml:"session,omitempty" json:"session,omitempty"`

		Cache struct {
			Enable   bool `yaml:"enable,omitempty" json:"enable,omitempty" default:"true"`
			UseRedis bool `yaml:"use_redis,omitempty" json:"use_redis,omitempty" default:"true"` // use redis cache
		} `yaml:"cache,omitempty" json:"cache,omitempty"`

		Static struct {
			Enable    bool   `yaml:"enable,omitempty" json:"enable,omitempty" default:"true"`
			StaticDir string `yaml:"static_dir,omitempty" json:"static_dir,omitempty" default:"./static"`
			StaticURL string `yaml:"static_url,omitempty" json:"static_url,omitempty" default:"/static"`
			Indexes   bool   `yaml:"indexes,omitempty" json:"indexes,omitempty" default:"true"`
		} `yaml:"static,omitempty" json:"static,omitempty"`

		Gzip struct {
			Enable bool `yaml:"enable,omitempty" json:"enable,omitempty" default:"true"`
		} `yaml:"gzip,omitempty" json:"gzip,omitempty"`
	} `yaml:"middlewares,omitempty" json:"middlewares,omitempty"`
}
