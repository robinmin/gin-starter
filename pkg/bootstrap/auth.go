package bootstrap

import (
	"context"
	"database/sql"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	cadapter "github.com/memwey/casbin-sqlx-adapter"
	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/robinmin/gin-starter/pkg/internal/dbo"
	_ "github.com/robinmin/gin-starter/pkg/internal/dbo"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 用户结构
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Email        string
}

// 角色结构
type Role struct {
	ID          int
	Name        string
	Description string
}

// 权限结构
type Permission struct {
	ID          int
	Name        string
	Description string
}

// 认证器
type Authenticator struct {
	db *DBToolKit
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// NewAuthenticator 创建认证器
func NewAuthenticator(db *DBToolKit) *Authenticator {
	return &Authenticator{
		db: db,
	}
}

// Authenticate 用户认证
func (a *Authenticator) Authenticate(username, password string) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	q := dbo.New(a.db)
	_, err = q.GetValidUserInfo(context.Background(), username, string(hashedPassword))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	// TODO: cache user information
	return true, nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 授权器
type Authorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(cfg types.AppConfig, logger *AppLogger) (*Authorizer, error) {
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
	enfcer, err := casbin.NewEnforcer(cfg.Middlewares.Auth.ModelFile, cadapter.NewAdapterFromOptions(opts))
	if err != nil {
		logger.Error("Failed to get DB connection string: " + err.Error())
		return nil, err
	}
	return &Authorizer{enforcer: enfcer}, err
}

func NewAuthorizerWithDB(model_file string, dbkit *DBToolKit) (*Authorizer, error) {
	opts := &cadapter.AdapterOptions{
		DB: (*sqlx.DB)(dbkit),
	}

	// Casbin v2 may return an error
	enfcer, err := casbin.NewEnforcer(model_file, cadapter.NewAdapterFromOptions(opts))
	return &Authorizer{enforcer: enfcer}, err
}

// HasPermission 检查用户是否拥有权限
func (author *Authorizer) HasPermission(user string, permission string) bool {
	result, err := author.enforcer.Enforce(user, permission, "*")
	if err != nil {
		return false
	}
	return result
}

func (author *Authorizer) AuthorizerHandler() gin.HandlerFunc {
	if author.enforcer == nil {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	return authz.NewAuthorizer(author.enforcer)
}
