package middleware

import (
	"github.com/gin-gonic/gin"
)

func SkippedPathPrefixes(ctx *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return false
	}

	path := ctx.Request.URL.Path
	pathLen := len(path)
	for _, prefix := range prefixes {
		if pl := len(prefix); pathLen >= pl && path[:pl] == prefix {
			return true
		}
	}
	return false
}

func AllowedPathPrefixes(ctx *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return true
	}

	path := ctx.Request.URL.Path
	pathLen := len(path)
	for _, prefix := range prefixes {
		if pl := len(prefix); pathLen >= pl && path[:pl] == prefix {
			return true
		}
	}
	return false
}

func Empty() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
