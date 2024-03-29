package bootstrap

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Module("bootstrap",
	fx.Provide(
		NewAppLogger,
		NewRedisPool,
		NewDB,
		NewSentry,
		NewApplication,
		NewCache,
		// NewHttpServer,
	),
)
