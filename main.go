package main

import (
	"StockPaperTradingApp/db"
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

	server.Run("localhost:8080")
}
