package routes

import (
	"StockPaperTradingApp/controllers"

	"github.com/gin-gonic/gin"
)

var (
	userController controllers.UserController = controllers.New()
)

// Used to create a user.
// Expects a first name, last name, email, and password. The email must not be already exist.
// Returns an auth token and user's first name, last name, and email
func RegisterEndpoint(ctx *gin.Context) {
	ctx.JSON(userController.Register(ctx))
}

// Used get user information (login)
// Expects a valid email and password.
// Returns an auth token and user's first name, last name, and email
func LoginEnpdpoint(ctx *gin.Context) {
	ctx.JSON(userController.Login(ctx))
}