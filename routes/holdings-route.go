package routes

import (
	"StockPaperTradingApp/controllers"

	"github.com/gin-gonic/gin"
)

var (
	holdingsController controllers.HoldingsController = controllers.Holdings()
)

func CreateHoldingsEndpoint(ctx *gin.Context) {
	ctx.JSON(holdingsController.CreateHoldings(ctx))
}

func GetAllHoldingsEndpoint(ctx *gin.Context) {
	ctx.JSON(holdingsController.GetAllHoldings(ctx))
}

func GetHoldingsEndpoint(ctx *gin.Context) {
	ctx.JSON(holdingsController.GetHoldings(ctx))
}

func UpdateHoldingsEndpoint(ctx *gin.Context) {
	ctx.JSON(holdingsController.UpdateQuantityHoldings(ctx))
}

func DeleteHoldingsEndpoint(ctx *gin.Context) {
	ctx.JSON(holdingsController.DeleteHoldings(ctx))
}
