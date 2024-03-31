package controllers

import (
	"StockPaperTradingApp/db"
	"StockPaperTradingApp/models"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// methods
type ApiController interface {
	BuyStock(ctx *gin.Context) (int, gin.H)
	SellStock(ctx *gin.Context) (int, gin.H)
	GetAllData(ctx *gin.Context) (int, gin.H)
}

// varables
type apiController struct {
	helper HelperController
}

// contructor
func Api() ApiController {
	return &apiController{
		helper: Helper(),
	}
}

// symbol, quant
func (c *apiController) BuyStock(ctx *gin.Context) (int, gin.H) {
	userIdHex, _ := ctx.Get("user_id")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	// get body
	var requestBody RequestBody
	ctx.BindJSON(&requestBody)
	if requestBody.Quantity == 0 {
		return http.StatusBadRequest, gin.H{
			"message": "Must give quantity amount",
		}
	}

	// get price, companyName
	url := "https://yfapi.net/v6/finance/quote?region=US&lang=en&symbols=" + requestBody.Symbol
	stockInfo := c.helper.SendRequest(url)
	if len(stockInfo["quoteResponse"].(map[string]any)["result"].([]any)) == 0 {
		return http.StatusBadRequest, gin.H{
			"message": "stock does not exits",
		}
	}
	price := stockInfo["quoteResponse"].(map[string]any)["result"].([]any)[0].(map[string]any)["regularMarketPrice"].(float64)
	companyName := stockInfo["quoteResponse"].(map[string]any)["result"].([]any)[0].(map[string]any)["longName"].(string)

	// update user cash
	userFilter := bson.D{{Key: "_id", Value: userId}}
	var resUser models.User
	db.GetUserCollection().FindOne(context.TODO(), userFilter).Decode(&resUser)
	if resUser.Cash < price*float64(requestBody.Quantity) {
		return http.StatusBadRequest, gin.H{
			"message": "Not enough cash",
		}
	}
	updateCash := resUser.Cash - price*float64(requestBody.Quantity)
	updateUser := bson.D{{Key: "$set", Value: bson.D{{Key: "Cash", Value: updateCash}}}}
	db.GetUserCollection().UpdateOne(context.TODO(), userFilter, updateUser)

	// create activity
	activity := models.Activity{
		Symbol:       requestBody.Symbol,
		CompanyName:  companyName,
		Quantity:     requestBody.Quantity,
		Side:         "BUY",
		Price:        price,
		Initiated_on: primitive.NewDateTimeFromTime(time.Now().UTC()),
		User_id:      userId,
	}
	db.GetActivityCollection().InsertOne(context.TODO(), activity)

	// check if person owns stock
	var holding models.Holdings
	holdingsFilter := bson.D{{Key: "user_id", Value: userId}, {Key: "symbol", Value: requestBody.Symbol}}
	db.GetHoldingsCollection().FindOne(context.TODO(), holdingsFilter).Decode(&holding)
	// if holding exist, update holdings
	fmt.Println(holding)
	if holding.Symbol == requestBody.Symbol {
		updateHolding := bson.D{{Key: "$set", Value: bson.D{{Key: "Quantity", Value: holding.Quantity + requestBody.Quantity}}}}
		db.GetHoldingsCollection().UpdateOne(context.TODO(), holdingsFilter, updateHolding)
	} else {
		// else, create holdings
		db.GetHoldingsCollection().InsertOne(context.TODO(), models.Holdings{
			CompanyName: companyName,
			Symbol:      requestBody.Symbol,
			Quantity:    requestBody.Quantity,
			User_id:     userId,
		})
	}

	// Newcash symbol holdingAmount
	return 200, gin.H{
		"message": "successfully bought",
	}
}

func (c *apiController) SellStock(ctx *gin.Context) (int, gin.H) {
	userIdHex, _ := ctx.Get("user_id")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))

	// get body
	var requestBody RequestBody
	ctx.BindJSON(&requestBody)
	if requestBody.Quantity == 0 {
		return http.StatusBadRequest, gin.H{
			"message": "Must give quantity amount",
		}
	}
	fmt.Println(userId)
	// check if person has stock
	var holding models.Holdings
	holdingsFilter := bson.D{{Key: "user_id", Value: userId}, {Key: "symbol", Value: requestBody.Symbol}}
	db.GetHoldingsCollection().FindOne(context.TODO(), holdingsFilter).Decode(&holding)
	fmt.Println(holding)
	if holding.CompanyName == "" {
		return http.StatusBadRequest, gin.H{
			"message": "user does not own stock",
		}
	}
	// check if quant less than equal to number there holding
	if holding.Quantity < requestBody.Quantity {
		return http.StatusBadRequest, gin.H{
			"message": "user does not own enough number of stocks",
		}
	}
	// update or delete holding
	if holding.Quantity == requestBody.Quantity {
		//delete holding
		db.GetHoldingsCollection().DeleteOne(context.TODO(), holdingsFilter)
	} else {
		// update
		updateHolding := bson.D{{Key: "$set", Value: bson.D{{Key: "Quantity", Value: holding.Quantity - requestBody.Quantity}}}}
		db.GetHoldingsCollection().UpdateOne(context.TODO(), holdingsFilter, updateHolding)
	}
	// get stock price
	url := "https://yfapi.net/v6/finance/quote?region=US&lang=en&symbols=" + requestBody.Symbol
	stockInfo := c.helper.SendRequest(url)
	if len(stockInfo["quoteResponse"].(map[string]any)["result"].([]any)) == 0 {
		return http.StatusBadRequest, gin.H{
			"message": "stock does not exits",
		}
	}
	price := stockInfo["quoteResponse"].(map[string]any)["result"].([]any)[0].(map[string]any)["regularMarketPrice"].(float64)
	companyName := stockInfo["quoteResponse"].(map[string]any)["result"].([]any)[0].(map[string]any)["longName"].(string)

	// update user cash
	userFilter := bson.D{{Key: "_id", Value: userId}}
	var resUser models.User
	db.GetUserCollection().FindOne(context.TODO(), userFilter).Decode(&resUser)
	if resUser.Cash < price*float64(requestBody.Quantity) {
		return http.StatusBadRequest, gin.H{
			"message": "Not enough cash",
		}
	}
	updateCash := resUser.Cash + price*float64(requestBody.Quantity)
	updateUser := bson.D{{Key: "$set", Value: bson.D{{Key: "Cash", Value: updateCash}}}}
	db.GetUserCollection().UpdateOne(context.TODO(), userFilter, updateUser)

	// create activity
	activity := models.Activity{
		Symbol:       requestBody.Symbol,
		CompanyName:  companyName,
		Quantity:     requestBody.Quantity,
		Side:         "SELL",
		Price:        price,
		Initiated_on: primitive.NewDateTimeFromTime(time.Now().UTC()),
		User_id:      userId,
	}
	db.GetActivityCollection().InsertOne(context.TODO(), activity)
	return 200, gin.H{}
}

// dashboard, trending, acitivity, holdings
func (c *apiController) GetAllData(ctx *gin.Context) (int, gin.H) {
	userIdHex, _ := ctx.Get("user_id")
	userId, _ := primitive.ObjectIDFromHex(userIdHex.(string))
	userilter := bson.D{{Key: "user_id", Value: userId}}

	// Get all holdings
	cursorHoldings, _ := db.GetHoldingsCollection().Find(context.TODO(), userilter, options.Find())
	var holdings []models.Holdings
	cursorHoldings.All(context.TODO(), &holdings)

	// Get all activity
	cursorActivity, _ := db.GetActivityCollection().Find(context.TODO(), userilter, options.Find())
	var activities []models.Activity
	cursorActivity.All(context.TODO(), &activities)

	// Get trending
	trendingUrl := baseURL + "/v1/finance/trending/US"
	var res = c.helper.SendRequest(trendingUrl)
	listOfTrending := res["finance"].(map[string]any)["result"].([]any)[0].(map[string]any)["quotes"].([]any)

	listOfTrendingSymbols := []string{}
	for _, symbol := range listOfTrending {
		listOfTrendingSymbols = append(listOfTrendingSymbols, symbol.(map[string]any)["symbol"].(string))
	}
	trendingResults := c.helper.GetStockInformation(listOfTrendingSymbols)

	// **** Get dashboard ****
	encodedUrl := baseURL + "/v8/finance/spark?interval=1d&range=3mo&" + url.PathEscape("symbols=^GSPC")
	snpRes := c.helper.SendRequest(encodedUrl)

	// get networth
	filter := bson.D{{Key: "user_id", Value: userId}}
	cursor, _ := db.GetNetworthCollection().Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "initiated_on", Value: -1}}))
	var netWorthModelList []models.Networth
	cursor.All(context.TODO(), &netWorthModelList)
	var networthList []float64
	for _, h := range netWorthModelList {
		networthList = append(networthList, h.Networth)
	}

	// Get holdings and turn to symbol : quant
	var listOfSymbols []string
	symbolToQuantity := make(map[string]int)
	allHoldings := c.helper.GetHoldings(userId)
	for _, h := range allHoldings {
		symbolToQuantity[h.Symbol] = h.Quantity
		listOfSymbols = append(listOfSymbols, h.Symbol)
	}

	// calculate asset worth
	symbolsInformation := c.helper.GetStockInformation(listOfSymbols)
	symbolToWorth := make(map[string]float64)
	var assetsWorth = 0.00
	for _, symbolInfo := range symbolsInformation {
		symbol := symbolInfo.(map[string]any)["symbol"].(string)
		price := symbolInfo.(map[string]any)["regularMarketPrice"].(float64)
		quantity := symbolToQuantity[symbol]
		assetWorth := price * float64(quantity)
		//for next section sector info
		symbolToWorth[symbol] = assetWorth
		assetsWorth += assetWorth
	}

	// get secotor info
	SectorToPercentage := make(map[string]float64)
	for _, symbol := range listOfSymbols {
		rawUrl := baseURL + "/v11/finance/quoteSummary/" + symbol + "?lang=en&region=US&modules=assetProfile"
		result := c.helper.SendRequest(rawUrl)
		sector := result["quoteSummary"].(map[string]any)["result"].([]any)[0].(map[string]any)["assetProfile"].(map[string]any)["sector"].(string)
		SectorToPercentage[sector] += symbolToWorth[symbol]
	}

	return http.StatusOK, gin.H{
		"holdings":   holdings,
		"activities": activities,
		"trending":   trendingResults,
		"dashboard": gin.H{
			"assetsWorth":    assetsWorth,
			"diversityGraph": SectorToPercentage,
			"performaceGraph": gin.H{
				"timeStamp":    snpRes["^GSPC"].(map[string]any)["timestamp"],
				"snpPrice":     snpRes["^GSPC"].(map[string]any)["close"],
				"netWorthList": networthList,
			},
		},
	}
}

type RequestBody struct {
	Quantity int
	Symbol   string
}
