package middleware

import (
	"fmt"
	"strings"

	// "github.com/LyricTian/gin-admin/v10/pkg/logging"
	// "github.com/LyricTian/gin-admin/v10/pkg/utility"

	"github.com/robinmin/gin-starter/pkg/utility"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type TraceConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	RequestHeaderKey    string
	ResponseTraceKey    string
}

var DefaultTraceConfig = TraceConfig{
	RequestHeaderKey: "X-Request-Id",
	ResponseTraceKey: "X-Trace-Id",
}

func Trace() gin.HandlerFunc {
	return TraceWithConfig(DefaultTraceConfig)
}

func TraceWithConfig(config TraceConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !AllowedPathPrefixes(ctx, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(ctx, config.SkippedPathPrefixes...) {
			ctx.Next()
			return
		}

		traceID := ctx.GetHeader(config.RequestHeaderKey)
		if traceID == "" {
			traceID = fmt.Sprintf("TRACE-%s", strings.ToUpper(xid.New().String()))
		}

		_ctx := utility.NewTraceID(ctx.Request.Context(), traceID)
		// TODO: enable this line once everything is ready
		// _ctx = logging.NewTraceID(_ctx, traceID)
		ctx.Request = ctx.Request.WithContext(_ctx)
		ctx.Writer.Header().Set(config.ResponseTraceKey, traceID)
		ctx.Next()
	}
}
