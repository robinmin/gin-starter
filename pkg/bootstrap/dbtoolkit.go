package bootstrap

import (
	"fmt"

	// _ "gorm.io/driver/sqlite" // // Sqlite driver based on GGO
	_ "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/jmoiron/sqlx"
)

type DBParams struct {
	// fx.In

	Type     string
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

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

type DBToolKit struct {
	*sqlx.DB
}

func NewDB(param DBParams) (*DBToolKit, error) {
	conn_str, err := param.GetConnectionString()
	if err != nil {
		return nil, fmt.Errorf("unsupported database type: %s", param.Type)
	}

	db, err1 := sqlx.Connect(param.Type, conn_str)
	return &DBToolKit{DB: db}, err1
}

// func NewQuery[T any](db *DBToolKit) *T {
// 	// return [T]{db: db}
// 	return new(T)
// }
