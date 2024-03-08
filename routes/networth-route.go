package routes

import (
	"StockPaperTradingApp/controllers"

	"github.com/gin-gonic/gin"
)

var (
	networthController controllers.NetworthController = controllers.Networth()
)

func GetAllNetworthEndpoint(ctx *gin.Context) {
	ctx.JSON(networthController.GetAllNetworth(ctx))
}
