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

type NetworthController interface {
	// CreateNetworth(ctx *gin.Context) (int, gin.H)
	GetAllNetworth(ctx *gin.Context) (int, gin.H)
}

type networthController struct{}

func Networth() NetworthController {
	return &networthController{}
}

// func (c *networthController) CreateNetworth(ctx *gin.Context) (int, gin.H) {}

func (c *networthController) GetAllNetworth(ctx *gin.Context) (int, gin.H) {
	res, ok := ctx.Get("user_id")
	if !ok {
		return http.StatusBadRequest, gin.H{
			"message": "Must have auth token",
		}
	}
	id, _ := primitive.ObjectIDFromHex(res.(string))
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, err := db.GetNetworthCollection().Find(context.TODO(), filter, options.Find())
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	var results []models.Networth
	if err = cursor.All(context.TODO(), &results); err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"networth": results,
	}
}
