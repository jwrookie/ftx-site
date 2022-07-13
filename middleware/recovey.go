package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/foxdex/ftx-site/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				log.Log.Error("recover", zap.Any("error:", err), zap.String("stack:", string(stack)))
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			}
		}()

		ctx.Next()
	}
}
