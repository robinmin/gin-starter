package config

import (
	"log/slog"

	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
)

const (
	// execute host command on local machine
	EVT_CLIENT_INIT types.UserDefinedEvent = iota
	EVT_CLIENT_CLOSE
)

var SentryEventsMeta = types.UserDefinedEventMap{
	EVT_CLIENT_INIT:  {Name: "evt_client_init", Level: slog.LevelInfo, Group: "sys"},
	EVT_CLIENT_CLOSE: {Name: "evt_client_close", Level: slog.LevelInfo, Group: "sys"},
}
