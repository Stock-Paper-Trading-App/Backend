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

type HoldingsController interface {
	CreateHoldings(ctx *gin.Context) (int, gin.H)
	GetAllHoldings(ctx *gin.Context) (int, gin.H)
	GetHoldings(ctx *gin.Context) (int, gin.H)
	UpdateQuantityHoldings(ctx *gin.Context) (int, gin.H)
	DeleteHoldings(ctx *gin.Context) (int, gin.H)
}

// varables
type holdingsController struct{}

// contructor
func Holdings() HoldingsController {
	return &holdingsController{}
}

func (c *holdingsController) CreateHoldings(ctx *gin.Context) (int, gin.H) {
	var holdings models.Holdings
	ctx.BindJSON(&holdings)
	if holdings.CompanyName == "" || holdings.Quantity == 0 || holdings.Symbol == "" {
		return http.StatusBadRequest, gin.H{
			"message": "All fields must be filled",
		}
	}
	res, ok := ctx.Get("user_id")
	if !ok {
		return http.StatusBadRequest, gin.H{
			"message": "Must have auth token",
		}
	}
	holdings.User_id, _ = primitive.ObjectIDFromHex(res.(string))

	result, err := db.GetHoldingsCollection().InsertOne(context.TODO(), holdings)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	holdings.ID = result.InsertedID.(primitive.ObjectID)
	return http.StatusCreated, gin.H{
		"holdings": holdings,
	}
}

func (c *holdingsController) GetAllHoldings(ctx *gin.Context) (int, gin.H) {
	res, ok := ctx.Get("user_id")
	if !ok {
		return http.StatusBadRequest, gin.H{
			"message": "Must have auth token",
		}
	}
	id, _ := primitive.ObjectIDFromHex(res.(string))
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, err := db.GetHoldingsCollection().Find(context.TODO(), filter, options.Find())
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	var results []models.Holdings
	if err = cursor.All(context.TODO(), &results); err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"holdings": results,
	}
}

func (c *holdingsController) GetHoldings(ctx *gin.Context) (int, gin.H) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	var holding models.Holdings
	err := db.GetHoldingsCollection().FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&holding)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"holding": holding,
	}
}

type QuantityRequestBody struct {
	Quantity int
}

func (c *holdingsController) UpdateQuantityHoldings(ctx *gin.Context) (int, gin.H) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	var quantityRequestBody QuantityRequestBody
	ctx.BindJSON(&quantityRequestBody)
	if quantityRequestBody.Quantity == 0 {
		return http.StatusBadRequest, gin.H{
			"message": "Must give quantity amount",
		}
	}
	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "quantity", Value: quantityRequestBody.Quantity}}}}
	_, err := db.GetHoldingsCollection().UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	var holding models.Holdings
	e := db.GetHoldingsCollection().FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&holding)
	if e != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	return http.StatusOK, gin.H{
		"updatedHolding": holding,
	}
}

func (c *holdingsController) DeleteHoldings(ctx *gin.Context) (int, gin.H) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	result, err := db.GetHoldingsCollection().DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return http.StatusInternalServerError, gin.H{
			"message": "Something went wrong connecting to database",
			"error":   err,
		}
	}
	if result.DeletedCount == 1 {
		return http.StatusOK, gin.H{
			"message": "Deleted successfully",
		}
	}
	return http.StatusOK, gin.H{
		"message": "Did not delete",
	}
}
