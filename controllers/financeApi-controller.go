package controllers

import (
	"StockPaperTradingApp/db"
	"StockPaperTradingApp/models"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
type financeController struct{}

// contructor
func FinanceApi() FinanceController {
	return &financeController{}
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
	var res = SendRequest(url)
	return http.StatusOK, res
}

func (c *financeController) GetTrending(ctx *gin.Context) (int, gin.H) {
	url := baseURL + "/v1/finance/trending/US"
	var res = SendRequest(url)

	listOfTrending := res["finance"].(map[string]any)["result"].([]any)[0].(map[string]any)["quotes"].([]any)

	listOfTrendingSymbols := []string{}
	for _, symbol := range listOfTrending {
		listOfTrendingSymbols = append(listOfTrendingSymbols, symbol.(map[string]any)["symbol"].(string))
	}
	results := GetStockInformation(listOfTrendingSymbols)
	return http.StatusOK, gin.H{
		"res": results,
	}
}

func (c *financeController) GetDashboardInformation(ctx *gin.Context) (int, gin.H) {
	// snp performace - sector info - assets worth (holdings stock prices * holding quantity)  - networth
	res, _ := ctx.Get("user_id")
	id, _ := primitive.ObjectIDFromHex(res.(string))

	// snp
	encodedUrl := baseURL + "/v8/finance/spark?interval=1d&range=3mo&" + url.PathEscape("symbols=^GSPC")
	snpRes := SendRequest(encodedUrl)

	// get networth
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, _ := db.GetNetworthCollection().Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "initiated_on", Value: -1}}))
	var netWorthList []models.Networth
	cursor.All(context.TODO(), &netWorthList)

	// Get holdings and turn to symbol : quant
	var listOfSymbols []string
	symbolToQuantity := make(map[string]int)
	holdings := GetHoldings(id)
	for _, h := range holdings {
		symbolToQuantity[h.Symbol] = h.Quantity
		listOfSymbols = append(listOfSymbols, h.Symbol)
	}

	// calculate asset worth
	symbolsInformation := GetStockInformation(listOfSymbols)
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
		result := SendRequest(rawUrl)
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

// Helpers
func GetStockInformation(symbols []string) []any {
	var quriesList [][]string
	var curr []string
	for _, sym := range symbols {
		curr = append(curr, sym)
		if len(curr) == 10 {
			quriesList = append(quriesList, curr)
			curr = []string{}
		}
	}
	if len(curr) > 0 {
		quriesList = append(quriesList, curr)
	}

	var results []any
	rawURL := baseURL + "/v6/finance/quote?region=US&lang=en&symbols="
	for _, queryArray := range quriesList {
		query := strings.Join(queryArray, ",")
		encodedquery := url.PathEscape(query)
		rawURL += encodedquery
		var res = SendRequest(rawURL)
		information := res["quoteResponse"].(map[string]any)["result"].([]any)
		results = append(results, information...)
	}
	return results
}

func SendRequest(rawUrl string) map[string]any {
	var apiToken = os.Getenv("API_KEY")

	r, _ := http.NewRequest(http.MethodGet, rawUrl, nil)

	r.Header.Set("x-api-key", apiToken)
	client := &http.Client{}
	resp, _ := client.Do(r)
	// read response body
	body, _ := io.ReadAll(resp.Body)

	// close response body (idk)
	defer resp.Body.Close()

	var res map[string]any
	json.Unmarshal(body, &res)
	return res
}

func GetHoldings(id primitive.ObjectID) []models.Holdings {
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, _ := db.GetHoldingsCollection().Find(context.TODO(), filter, options.Find())
	var results []models.Holdings
	cursor.All(context.TODO(), &results)
	return results
}

func UpdateNetworths() {
	// Get all users
	cursor, _ := db.GetUserCollection().Find(context.TODO(), bson.D{}, options.Find())
	var users []models.User
	cursor.All(context.TODO(), &users)

	// for each user
	for _, user := range users {
		// get holdings
		holdings := GetHoldings(user.ID)
		// calculate asset
		var listOfSymbols []string
		symbolToQuantity := make(map[string]int)
		for _, h := range holdings {
			symbolToQuantity[h.Symbol] = h.Quantity
			listOfSymbols = append(listOfSymbols, h.Symbol)
		}

		// calculate asset worth
		symbolsInformation := GetStockInformation(listOfSymbols)
		var assetsWorth = 0.00
		for _, symbolInfo := range symbolsInformation {
			symbol := symbolInfo.(map[string]any)["symbol"].(string)
			price := symbolInfo.(map[string]any)["regularMarketPrice"].(float64)
			quantity := symbolToQuantity[symbol]
			assetsWorth += price * float64(quantity)
		}

		// save networth
		var networth = models.Networth{
			Networth:     user.Cash + int(assetsWorth),
			Initiated_on: primitive.NewDateTimeFromTime(time.Now().UTC()),
			User_id:      user.ID,
		}
		db.GetNetworthCollection().InsertOne(context.TODO(), networth)
	}
}
