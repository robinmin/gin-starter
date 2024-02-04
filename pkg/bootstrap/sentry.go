package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"

	"github.com/robinmin/gin-starter/pkg/bootstrap/types"
)

type AppSentry struct {
	Params types.AppSentryConfig
}

func NewSentry(cfg types.AppConfig, lc fx.Lifecycle, logger *AppLogger) (*AppSentry, error) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := sentry.Init(sentry.ClientOptions{
				Dsn: cfg.Sentry.DSN,
				// Set TracesSampleRate to 1.0 to capture 100%
				// of transactions for performance monitoring.
				// We recommend adjusting this value in production,
				TracesSampleRate: cfg.Sentry.TracesSampleRate,
			})
			if err != nil {
				logger.Error("Sentry init error : " + err.Error())
				return err
			}

			logger.Info("Sentry started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// RecoverPanic 恢复panic并发送到sentry
			// 兜底所有的异常监控与处理
			err := recover()
			if err != nil {
				sentry.CurrentHub().Recover(err)

				// 确保所有事件都被发送到Sentry
				sentry.Flush(time.Second * 5)
			}

			logger.Info("Sentry stopped")
			return nil
		},
	})
	return &AppSentry{Params: cfg.Sentry}, nil
}

func (*AppSentry) SetUser(id string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID: id,
		})
	})
}

func (*AppSentry) SetTag(key, value string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

func (*AppSentry) SetExtra(key string, value interface{}) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra(key, value)
	})
}

// CaptureException 捕获异常并发送到sentry
func (*AppSentry) CaptureException(err error) {
	sentry.CaptureException(err)
}

// CaptureRequest 捕获请求并发送到sentry
func (*AppSentry) CaptureRequest(r *http.Request) {
	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetRequest(r)
}

// // ReportCustomEvent 上报定制事件
func (sty *AppSentry) ReportEvent(event_id types.UserDefinedEvent, eventMessage string, payLoad map[string]interface{}) {
	needReport, meta := sty.getEventConfig(event_id)
	if needReport {
		event := sentry.NewEvent()
		event.Level = sentry.Level(meta.Level)
		event.Message = eventMessage
		event.Tags = map[string]string{
			"event_name": meta.Name,
			"event_id":   uuid.New().String(),
		}
		if payLoad != nil {
			event.Extra = payLoad
		}
		sentry.CaptureEvent(event)
	}
}

func (sty *AppSentry) getEventConfig(event_id types.UserDefinedEvent) (bool, *types.UserDefinedEventMeta) {
	var meta types.UserDefinedEventMeta
	var ok bool
	var needReport bool

	if meta, ok = sty.Params.EventsMeta[event_id]; !ok {
		// by default report all
		meta = types.UserDefinedEventMeta{
			Name:  "evnt_unknown_report",
			Level: "debug",
			Group: "sys",
		}
	}

	level, err := zapcore.ParseLevel(meta.Level)
	if err == nil && int(level) >= sty.Params.DefaultLevel {
		needReport = true
	} else {
		needReport = false
	}

	return needReport, &meta
}
