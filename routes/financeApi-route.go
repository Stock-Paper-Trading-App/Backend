package routes

import (
	"StockPaperTradingApp/controllers"

	"github.com/gin-gonic/gin"
)

var (
	financeController controllers.FinanceController = controllers.FinanceApi()
)

func AutoCompleteEndpoint(ctx *gin.Context) {
	ctx.JSON(financeController.GetAutoComplete(ctx))
}

func TrendingEndpoint(ctx *gin.Context) {
	ctx.JSON(financeController.GetTrending(ctx))
}

func DashBoardEndpoint(ctx *gin.Context) {
	ctx.JSON(financeController.GetDashboardInformation(ctx))
}

func StockPageInformationEndpoint(ctx *gin.Context) {
	ctx.JSON(financeController.GetStockPageInformation(ctx))
}

func StockInformationEndpoint(ctx *gin.Context) {
	ctx.JSON(financeController.GetStockInformation(ctx))
}
