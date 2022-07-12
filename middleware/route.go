package middleware

import (
	"github.com/foxdex/ftx-site/config"
	"github.com/foxdex/ftx-site/handler"
	"github.com/gin-gonic/gin"
)

func NewRoute(api *gin.Engine) {
	var (
		luckyDrawHandler handler.LuckyDrawHandler
	)
	conf := config.GetConfig()
	root := api.Group(conf.App.Model)

	{
		lucky := root.Group("/lucky", Csrf())
		lucky.POST("/token", luckyDrawHandler.CreateToken)
		lucky.POST("/draw", Ticket(), luckyDrawHandler.Draw)
		lucky.POST("/award", Ticket(), luckyDrawHandler.Award)
		lucky.GET("/:email", luckyDrawHandler.GetResult)
		lucky.GET("/jackpot", luckyDrawHandler.GetJackpot)
	}
}