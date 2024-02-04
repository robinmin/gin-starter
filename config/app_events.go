package config

import (
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
)

const (
	// execute host command on local machine
	EVT_CLIENT_INIT types.UserDefinedEvent = iota
	EVT_CLIENT_CLOSE
)

var SentryEventsMeta = types.UserDefinedEventMap{
	EVT_CLIENT_INIT:  types.UserDefinedEventMeta{Name: "evt_client_init", Level: "info", Group: "sys"},
	EVT_CLIENT_CLOSE: types.UserDefinedEventMeta{Name: "evt_client_close", Level: "info", Group: "sys"},
}
