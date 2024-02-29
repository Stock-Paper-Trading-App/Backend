package controllers

import (
	"StockPaperTradingApp/db"
	"StockPaperTradingApp/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// methods
type UserController interface {
	Register(ctx *gin.Context) (int, gin.H)
	Login(ctx *gin.Context) (int, gin.H)
}

// varables
type authController struct{}

// contructor
func Auth() UserController {
	return &authController{}
}

func (c *authController) Register(ctx *gin.Context) (int, gin.H) {
	var user models.User
	ctx.BindJSON(&user)
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		return http.StatusBadRequest, gin.H{
			"message": "All fields must be filled",
		}
	}
	// check if email already exists
	filter := bson.D{{Key: "email", Value: user.Email}}
	var foundUser models.User
	err := db.GetUserCollection().FindOne(context.TODO(), filter).Decode(&foundUser)
	if err == nil {
		return http.StatusConflict, gin.H{
			"message": "Email already exists",
		}
	}
	// Save user
	user.HashPassword()
	user.Cash = 50000
	result, err := db.GetUserCollection().InsertOne(context.TODO(), user)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	token, err := user.CreateJWT()
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong creating JWT token",
			"error":   err,
		}
	}
	return http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"firstname": user.FirstName,
			"lastname":  user.LastName,
			"email":     user.Email,
			"cash":      user.Cash,
		},
	}
}

func (c *authController) Login(ctx *gin.Context) (int, gin.H) {
	var user models.User
	ctx.BindJSON(&user)
	if user.Email == "" || user.Password == "" {
		return http.StatusBadRequest, gin.H{
			"message": "All fields must be filled",
		}
	}
	var resultUser models.User
	filter := bson.D{{Key: "email", Value: user.Email}}
	err := db.GetUserCollection().FindOne(context.TODO(), filter).Decode(&resultUser)
	if err != nil || !user.ComparePasswords([]byte(resultUser.Password)) {
		return http.StatusUnauthorized, gin.H{
			"message": "Incorrect Credentials",
		}
	}
	token, err := resultUser.CreateJWT()
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong creating JWT token",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"firstname": resultUser.FirstName,
			"lastname":  resultUser.LastName,
			"email":     resultUser.Email,
			"cash":      resultUser.Cash,
		},
	}
}
