package bootstrap

import (
	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
	cadapter "github.com/memwey/casbin-sqlx-adapter"
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
)

func NewAuthEnforcer(cfg types.AppConfig, logger *AppLogger) (*casbin.Enforcer, error) {
	param := DBParams(types.AppDBConfig{
		Type:     cfg.Database.Type,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Database: cfg.Database.Database,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	})
	connection_str, err := param.GetDSN()
	if err != nil {
		logger.Error("Failed to get DB connection string: " + err.Error())
		return nil, err
	}

	opts := &cadapter.AdapterOptions{
		DriverName:     cfg.Database.Type,
		DataSourceName: connection_str,
		TableName:      cfg.Middlewares.Auth.TableName,
		// or reuse an existing connection:
		// DB: myDBConn,
	}

	// Casbin v2 may return an error
	return casbin.NewEnforcer(cfg.Middlewares.Auth.ModelFile, cadapter.NewAdapterFromOptions(opts))
}

func NewAuthEnforcerFromDB(model_file string, dbkit *DBToolKit) (*casbin.Enforcer, error) {
	opts := &cadapter.AdapterOptions{
		DB: (*sqlx.DB)(dbkit),
	}

	// Casbin v2 may return an error
	return casbin.NewEnforcer(model_file, cadapter.NewAdapterFromOptions(opts))
}
