package bootstrap

import (
	"fmt"

	// _ "gorm.io/driver/sqlite" // // Sqlite driver based on GGO
	_ "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/jmoiron/sqlx"
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
)

type DBParams types.AppDBConfig
type DBToolKit sqlx.DB

func NewDBParams(dbtype string, dbhost string, dbport int, dbdatabase string, dbuser string, dbpassword string) DBParams {
	return DBParams{
		Type:     dbtype,
		Host:     dbhost,
		Port:     dbport,
		Database: dbdatabase,
		User:     dbuser,
		Password: dbpassword,
	}
}

func (param DBParams) GetConnectionString() (string, error) {
	switch param.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", param.User, param.Password, param.Host, param.Port, param.Database), nil
	case "sqlite3":
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", param.Database), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", param.Type)
	}
}

func NewDB(cfg types.AppConfig) (*DBToolKit, error) {
	params := DBParams(cfg.Database)
	conn_str, err0 := params.GetConnectionString()
	if err0 != nil {
		return nil, fmt.Errorf("Unsupported database type: %s", params.Type)
	}

	db, err := sqlx.Connect(params.Type, conn_str)
	return (*DBToolKit)(db), err
}
