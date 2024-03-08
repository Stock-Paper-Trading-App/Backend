package main

import (
	"StockPaperTradingApp/db"
	"StockPaperTradingApp/middlewares"
	"StockPaperTradingApp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	db.ConnectToDB()

	auth := server.Group("/auth")
	{
		auth.POST("/register", routes.RegisterEndpoint)
		auth.GET("/login", routes.LoginEnpdpoint)
	}

	holdings := server.Group("/holdings").Use(middlewares.Authentication)
	{
		holdings.POST("/", routes.CreateHoldingsEndpoint)
		holdings.GET("/", routes.GetAllHoldingsEndpoint)
		holdings.GET("/:id", routes.GetHoldingsEndpoint)
		holdings.PATCH("/:id", routes.UpdateHoldingsEndpoint)
		holdings.DELETE("/:id", routes.DeleteHoldingsEndpoint)
	}

	activity := server.Group("/activity").Use(middlewares.Authentication)
	{
		activity.POST("/", routes.CreateHoldingsEndpoint)
		activity.GET("/", routes.GetAllHoldingsEndpoint)
		activity.GET("/:id", routes.GetHoldingsEndpoint)
	}

	server.Run("localhost:8080")
}
