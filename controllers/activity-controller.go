package controllers

import (
	"StockPaperTradingApp/db"
	"StockPaperTradingApp/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActivityController interface {
	CreateActivity(ctx *gin.Context) (int, gin.H)
	GetAllActivity(ctx *gin.Context) (int, gin.H)
	GetActivity(ctx *gin.Context) (int, gin.H)
	// UpdateQuantityActivity(ctx *gin.Context) (int, gin.H)
	// DeleteActivity(ctx *gin.Context) (int, gin.H)
}

// varables
type activityController struct{}

// contructor
func Activity() ActivityController {
	return &activityController{}
}

func (c *activityController) CreateActivity(ctx *gin.Context) (int, gin.H) {
	var activity models.Activity
	ctx.BindJSON(&activity)

	if activity.CompanyName == "" || activity.Quantity == 0 || activity.Symbol == "" || activity.Side == "" || activity.Price == 0 || activity.Description == "" {
		return http.StatusBadRequest, gin.H{
			"message": "All fields must be filled",
		}
	}
	if activity.Side != "BUY" && activity.Side != "SELL" {
		return http.StatusBadRequest, gin.H{
			"message": "Side must be either BUY or SELL",
		}
	}
	res, ok := ctx.Get("user_id")
	if !ok {
		return http.StatusBadRequest, gin.H{
			"message": "Must have auth token",
		}
	}
	activity.User_id, _ = primitive.ObjectIDFromHex(res.(string))

	result, err := db.GetActivityCollection().InsertOne(context.TODO(), activity)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	activity.ID = result.InsertedID.(primitive.ObjectID)
	return http.StatusCreated, gin.H{
		"activity": activity,
	}
}

func (c *activityController) GetAllActivity(ctx *gin.Context) (int, gin.H) {
	res, ok := ctx.Get("user_id")
	if !ok {
		return http.StatusBadRequest, gin.H{
			"message": "Must have auth token",
		}
	}
	id, _ := primitive.ObjectIDFromHex(res.(string))
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, err := db.GetActivityCollection().Find(context.TODO(), filter, options.Find())
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	var results []models.Activity
	if err = cursor.All(context.TODO(), &results); err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"activity": results,
	}
}

func (c *activityController) GetActivity(ctx *gin.Context) (int, gin.H) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	var activity models.Activity
	err := db.GetActivityCollection().FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&activity)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"activity": activity,
	}
}
