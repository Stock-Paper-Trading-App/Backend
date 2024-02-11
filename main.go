package main

import (
	"StockPaperTradingApp/controllers"
	"StockPaperTradingApp/db"

	"github.com/gin-gonic/gin"
)

var (
	userController controllers.UserController = controllers.New()
)

func main() {
	server := gin.Default()
	db.ConnectToDB()

	server.GET("/register", func(ctx *gin.Context) {
		ctx.JSON(userController.Register(ctx))
	})

	server.GET("/login", func(ctx *gin.Context) {
		ctx.JSON(userController.Login(ctx))
	})

	server.Run("localhost:8080")
}
