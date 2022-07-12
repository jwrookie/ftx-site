package middleware

import (
	"net/http"

	"github.com/foxdex/ftx-site/pkg/consts"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if gin.IsDebugging() {
			//动态允许Allow-Origin，解决"*"不能与Access-Control-Allow-Credentials为true共存的问题
			ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
			//拿到跨域请求中Header的字段，不能用"*"

			ctx.Header("Access-Control-Allow-Headers", consts.HeaderAUTHCSRFTOKEN)
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,PUT,DELETE")
			ctx.Header("Access-Control-Allow-Credentials", "true")

			// 能够让客户端的js读取到Header
			ctx.Header("Access-Control-Expose-Headers", consts.HeaderAUTHCSRFTOKEN)
			// 表明在xxx秒内，不需要再发送预检验请求，可以缓存该结果
			ctx.Header("Access-Control-Max-Age", "3600")
			// 如果method是OPTIONS，直接返回成功
			if ctx.Request.Method == http.MethodOptions {
				ctx.AbortWithStatusJSON(http.StatusOK, "Options Request!")
				return
			}
		}

		// 处理请求
		ctx.Next()
	}
}
