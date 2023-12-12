package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"os"
	"sync"

	sloggin "github.com/samber/slog-gin"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type Application[T any] struct {
	// Configuration
	Config *T

	// Engine instance
	engine *gin.Engine

	// DB instance
	// DB     *database.DB

	// log
	logFileName    string
	logFileHandler *os.File
	Logger         *slog.Logger

	wg sync.WaitGroup
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewApplication[T any](cfg *T, logFile string) (*Application[T], error) {
	app := &Application[T]{Config: cfg}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	writer := app.createLogWriter(logFile)
	if writer == nil {
		gin.DefaultWriter = os.Stdout
	} else {
		gin.DefaultWriter = writer
	}
	app.Logger = slog.New(slog.NewTextHandler(gin.DefaultWriter, opts)) //.With("gin_mode", gin.EnvGinMode)

	config := sloggin.Config{
		WithSpanID:  true,
		WithTraceID: true,
	}

	app.engine = gin.New()
	// The middleware will log all requests attributes.
	app.engine.Use(sloggin.NewWithConfig(app.Logger, config), gin.Recovery())
	// app.engine.ForwardedByClientIP = true
	// app.engine.SetTrustedProxies([]string{"127.0.0.1"})
	// app.engine.SetTrustedProxies(strings.Split(app.Config.System.TrustedProxies, ";"))

	if gin.IsDebugging() {
		gin.ForceConsoleColor()
	} else {
		// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
		gin.DisableConsoleColor()
	}

	// Example pong request.
	app.engine.GET("/pong", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return app, nil
}

func (app *Application[T]) createLogWriter(filename string) io.Writer {
	var writers []io.Writer
	if gin.IsDebugging() {
		writers = append(writers, os.Stdout)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current working directory:", cwd)

	if app.logFileHandler == nil {
		var err error
		app.logFileHandler, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Failed to open log file: %w", err)
			return nil
		}
		app.logFileName = filename
	}

	writers = append(writers, app.logFileHandler)
	return io.MultiWriter(writers...)
}

func (app *Application[T]) RunServer(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.engine,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	app.Logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	app.Logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))

	app.wg.Wait()
	return nil
}

func (app *Application[T]) Quit() {
	// close log file handler when application exits
	if app.logFileHandler != nil {
		err := app.logFileHandler.Close()
		if err != nil {
			fmt.Println("Failed to close log file: %w", err)
		}
		app.logFileHandler = nil
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

// LoadConfig 从指定的YAML文件中加载配置信息
func LoadConfig[T any](yamlFile string) (*T, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	var config T
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 将配置信息保存到指定的YAML文件中
func SaveConfig[T any](cfg *T, yamlFile string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(yamlFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// func (cfg *AppConfig) GetConnectionString() (string, error) {
// 	switch cfg.Database.Type {
// 	case "mysql":
// 		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database), nil
// 	case "sqlite3":
// 		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", cfg.Database.Database), nil
// 	default:
// 		return "", fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
// 	}
// }
