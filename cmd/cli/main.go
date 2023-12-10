package main

import (
	log "log/slog"
	"os"
	"runtime/debug"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	logger := log.New(log.NewTextHandler(os.Stdout, &log.HandlerOptions{Level: log.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *log.Logger) error {
	// var cfg config
	log.Debug("Server starting......")
	return nil
}
