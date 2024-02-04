package bootstrap

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Module("bootstrap",
	fx.Provide(
		NewAppLogger,
		// NewDBParams,
		NewDB,
		NewSentry,
		NewApplication,
		NewRedisCache,
		NewMemoryCache,
		// NewHttpServer,
	),
)
