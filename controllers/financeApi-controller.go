package controllers

import (
	"StockPaperTradingApp/db"
	"StockPaperTradingApp/models"
	"context"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FinanceController interface {
	GetAutoComplete(ctx *gin.Context) (int, gin.H)
	GetTrending(ctx *gin.Context) (int, gin.H)
	GetStockPageInformation(ctx *gin.Context) (int, gin.H)
	GetDashboardInformation(ctx *gin.Context) (int, gin.H)
}

// varables
type financeController struct {
	helper HelperController
}

// contructor
func FinanceApi() FinanceController {
	return &financeController{
		helper: Helper(),
	}
}

var baseURL = "https://yfapi.net"

func (c *financeController) GetAutoComplete(ctx *gin.Context) (int, gin.H) {
	var queries = ctx.Request.URL.Query()
	var query = queries["query"]
	if len(query) != 1 {
		return http.StatusBadRequest, gin.H{
			"error": "Must have one value assigned to query in parameter",
		}
	}

	url := baseURL + "/v6/finance/autocomplete?region=US&lang=en&query=" + query[0]
	var res = c.helper.SendRequest(url)
	return http.StatusOK, res
}

func (c *financeController) GetTrending(ctx *gin.Context) (int, gin.H) {
	url := baseURL + "/v1/finance/trending/US"
	var res = c.helper.SendRequest(url)

	listOfTrending := res["finance"].(map[string]any)["result"].([]any)[0].(map[string]any)["quotes"].([]any)

	listOfTrendingSymbols := []string{}
	for _, symbol := range listOfTrending {
		listOfTrendingSymbols = append(listOfTrendingSymbols, symbol.(map[string]any)["symbol"].(string))
	}
	results := c.helper.GetStockInformation(listOfTrendingSymbols)
	return http.StatusOK, gin.H{
		"res": results,
	}
}

// break down
func (c *financeController) GetDashboardInformation(ctx *gin.Context) (int, gin.H) {
	// snp performace - sector info - assets worth (holdings stock prices * holding quantity)  - networth
	res, _ := ctx.Get("user_id")
	id, _ := primitive.ObjectIDFromHex(res.(string))

	// snp
	encodedUrl := baseURL + "/v8/finance/spark?interval=1d&range=3mo&" + url.PathEscape("symbols=^GSPC")
	snpRes := c.helper.SendRequest(encodedUrl)

	// get networth
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, _ := db.GetNetworthCollection().Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "initiated_on", Value: -1}}))
	var netWorthList []models.Networth
	cursor.All(context.TODO(), &netWorthList)

	// Get holdings and turn to symbol : quant
	var listOfSymbols []string
	symbolToQuantity := make(map[string]int)
	holdings := c.helper.GetHoldings(id)
	for _, h := range holdings {
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
		percentage := symbolToWorth[symbol] / assetsWorth
		SectorToPercentage[sector] += (percentage * 100)
	}

	// return info
	return 200, gin.H{
		"assetsWorth":    assetsWorth,
		"diversityGraph": SectorToPercentage,
		"performaceGraph": gin.H{
			"timeStamp":    snpRes["^GSPC"].(map[string]any)["timestamp"],
			"snpPrice":     snpRes["^GSPC"].(map[string]any)["close"],
			"netWorthList": netWorthList,
		},
	}
}

func (c *financeController) GetStockPageInformation(ctx *gin.Context) (int, gin.H) {
	var queries = ctx.Request.URL.Query()
	var query = queries["query"]
	if len(query) != 1 {
		return http.StatusBadRequest, gin.H{
			"error": "Must have one value assigned to query in parameter",
		}
	}
	// finish when we decide what to include (SPTA-38)

	return 200, gin.H{}

}