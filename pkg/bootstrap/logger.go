package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Debug      bool   `yaml:"debug,omitempty" json:"debug,omitempty" default:"true"`
	Level      string `yaml:"level,omitempty" json:"level,omitempty" default:"info"`
	CallerSkip int    `yaml:"caller_skip,omitempty" json:"caller_skip,omitempty" default:"2"`
	File       struct {
		Enable     bool   `yaml:"enable,omitempty" json:"enable,omitempty" default:"false"`
		Path       string `yaml:"path,omitempty" json:"path,omitempty" default:"./log/app.log"`
		MaxSize    int    `yaml:"maxsize,omitempty" json:"maxsize,omitempty" default:"100"`
		MaxBackups int    `yaml:"maxbackups,omitempty" json:"maxbackups,omitempty" default:"5"`
	} `yaml:"file,omitempty" json:"file,omitempty"`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type AppLogger struct {
	*zap.Logger
}

func NewAppLogger() *AppLogger {
	return &AppLogger{Logger: zap.L()}
}

func (logger *AppLogger) Print(v ...interface{}) {
	logger.Info(fmt.Sprint(v...))
}

func (logger *AppLogger) Write(p []byte) (n int, err error) {
	logger.Error(string(p))
	return len(p), nil
}

func (logger *AppLogger) GetRawLogger() *log.Logger {
	return log.New(logger, "", 0)
}

func InitLogger() (func(), error) {
	return InitLoggerWithConfig(
		NewInstance[LoggerConfig](),
	)
}

// InitLoggerWithConfig initializes the global logger with the given config
func InitLoggerWithConfig(cfg *LoggerConfig) (func(), error) {
	var zconfig zap.Config
	var encoder zapcore.Encoder
	var encoderConfig zapcore.EncoderConfig

	if cfg.Debug {
		cfg.Level = "debug"
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // Plain text for debug

		zconfig = zap.NewDevelopmentConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig) // JSON for non-debug

		zconfig = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zconfig.Level.SetLevel(level)

	var (
		logger   *zap.Logger
		cleanFns []func()
	)

	if cfg.File.Enable {
		filename := cfg.File.Path
		_ = os.MkdirAll(filepath.Dir(filename), 0777)
		fileWriter := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			Compress:   false,
			LocalTime:  true,
		}

		cleanFns = append(cleanFns, func() {
			_ = fileWriter.Close()
		})

		var ws []zapcore.WriteSyncer
		ws = append(ws, zapcore.AddSync(fileWriter))
		if cfg.Debug {
			ws = append(ws, zapcore.AddSync(os.Stdout)) // Add stdout in debug mode
		}

		zc := zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(ws...),
			zconfig.Level,
		)
		logger = zap.New(zc)
	} else {
		ilogger, err := zconfig.Build()
		if err != nil {
			return nil, err
		}
		logger = ilogger
	}

	skip := cfg.CallerSkip
	if skip <= 0 {
		skip = 2
	}

	logger = logger.WithOptions(
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(skip),
	)

	zap.ReplaceGlobals(logger)
	return func() {
		for _, fn := range cleanFns {
			fn()
		}

		if err := zap.L().Sync(); err != nil {
			fmt.Printf("failed to sync zap logger: %s \n", err.Error())
		}
	}, nil
}
