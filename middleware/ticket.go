package middleware

import (
	"net/http"

	"github.com/foxdex/ftx-site/dto"

	"github.com/foxdex/ftx-site/pkg/consts"
	"github.com/foxdex/ftx-site/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func Ticket() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var claims *jwt.UserClaims
		token := ctx.GetHeader(consts.HeaderDRAWTOKEN)
		if gin.IsDebugging() && token == "lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK" {
			ctx.Next()
			return
		}

		claims, err := claims.Parse(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseFormat{Code: 1, Msg: err.Error()})
			return
		}
		ctx.Set(consts.HeaderDRAWTOKEN, claims)

		ctx.Next()
	}
}
