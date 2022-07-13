package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/foxdex/ftx-site/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			payloadCopy string
			payload     []byte
		)
		reqMethod := c.Request.Method // 请求方式

		if reqMethod != http.MethodGet {
			payload, _ = io.ReadAll(c.Request.Body)
			payloadCopy = string(payload)
			c.Request.Body = io.NopCloser(bytes.NewReader(payload))
		}

		startTime := time.Now() // 开始时间
		c.Next()                // 处理请求
		endTime := time.Now()   // 结束时间

		latencyTime := endTime.Sub(startTime) // 执行时间

		reqUri := c.Request.RequestURI  // 请求方式
		statusCode := c.Writer.Status() // 状态码
		clientIP := c.ClientIP()        // 请求真实IP

		log.Log.Info("",
			zap.String("method", reqMethod),
			zap.String("uri", reqUri),
			zap.Int("code", statusCode),
			zap.Duration("latency", latencyTime),
			zap.String("client_ip", clientIP),
			zap.String("payload", payloadCopy),
		)
	}
}
