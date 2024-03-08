package routes

import (
	"StockPaperTradingApp/controllers"

	"github.com/gin-gonic/gin"
)

var (
	activityController controllers.ActivityController = controllers.Activity()
)

func CreateActivityEndpoint(ctx *gin.Context) {
	ctx.JSON(activityController.CreateActivity(ctx))
}

func GetAllActivityEndpoint(ctx *gin.Context) {
	ctx.JSON(activityController.GetAllActivity(ctx))
}

func GetActivityEndpoint(ctx *gin.Context) {
	ctx.JSON(activityController.GetActivity(ctx))
}
