package routes

import (
	"StockPaperTradingApp/controllers"

	"github.com/gin-gonic/gin"
)

var (
	apiController controllers.ApiController = controllers.Api()
)

func BuyStockEndpoint(ctx *gin.Context) {
	ctx.JSON(apiController.BuyStock(ctx))
}

func SellStockEndpoint(ctx *gin.Context) {
	ctx.JSON(apiController.SellStock(ctx))
}

func GetAllDataEndpoint(ctx *gin.Context) {
	ctx.JSON(apiController.GetAllData(ctx))
}
