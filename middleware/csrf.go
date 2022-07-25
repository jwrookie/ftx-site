package middleware

import (
	"net/http"
	"strconv"

	"github.com/foxdex/ftx-site/config"
	"github.com/foxdex/ftx-site/pkg/utils"

	"github.com/foxdex/ftx-site/pkg/consts"
	"github.com/gin-gonic/gin"
)

func Csrf() gin.HandlerFunc {
	skipUrl := map[string]struct{}{}

	return func(ctx *gin.Context) {
		if _, ok := skipUrl[ctx.Request.RequestURI]; ok {
			ctx.Next()
			return
		}

		token := ctx.GetHeader(consts.HeaderAUTHCSRFTOKEN)
		if gin.IsDebugging() && token == "lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK" {
			ctx.Next()
			return
		}
		if len(token) < 1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, nil)
			return
		}

		decryptToken, err := utils.Base64AESCBCDecrypt(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, nil)
			return
		}

		tokenUint, _ := strconv.ParseUint(decryptToken, 10, 64)
		now := utils.UnixMilliNow()

		// 检查时间误差
		_interval := uint64(config.GetConfig().Csrf.Interval)
		if (now > tokenUint && now-tokenUint > _interval) || (tokenUint > now && tokenUint-now > _interval) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, nil)
			return
		}

		ctx.Next()
	}
}
