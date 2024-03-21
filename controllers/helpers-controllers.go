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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HelperController interface {
	SendRequest(rawUrl string) map[string]any
	GetStockInformation(symbols []string) []any
	GetHoldings(id primitive.ObjectID) []models.Holdings
	UpdateNetworths()
}

// varables
type helperController struct{}

// contructor
func Helper() HelperController {
	return &helperController{}
}

func (c *helperController) SendRequest(rawUrl string) map[string]any {
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

func (c *helperController) GetStockInformation(symbols []string) []any {
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
		var res = c.SendRequest(rawURL)
		information := res["quoteResponse"].(map[string]any)["result"].([]any)
		results = append(results, information...)
	}
	return results
}

func (c *helperController) GetHoldings(id primitive.ObjectID) []models.Holdings {
	filter := bson.D{{Key: "user_id", Value: id}}
	cursor, _ := db.GetHoldingsCollection().Find(context.TODO(), filter, options.Find())
	var results []models.Holdings
	cursor.All(context.TODO(), &results)
	return results
}

func (c *helperController) UpdateNetworths() {
	// Get all users
	cursor, _ := db.GetUserCollection().Find(context.TODO(), bson.D{}, options.Find())
	var users []models.User
	cursor.All(context.TODO(), &users)

	// for each user
	for _, user := range users {
		// get holdings
		holdings := c.GetHoldings(user.ID)
		// calculate asset
		var listOfSymbols []string
		symbolToQuantity := make(map[string]int)
		for _, h := range holdings {
			symbolToQuantity[h.Symbol] = h.Quantity
			listOfSymbols = append(listOfSymbols, h.Symbol)
		}

		// calculate asset worth
		symbolsInformation := c.GetStockInformation(listOfSymbols)
		var assetsWorth = 0.00
		for _, symbolInfo := range symbolsInformation {
			symbol := symbolInfo.(map[string]any)["symbol"].(string)
			price := symbolInfo.(map[string]any)["regularMarketPrice"].(float64)
			quantity := symbolToQuantity[symbol]
			assetsWorth += price * float64(quantity)
		}

		// save networth
		var networth = models.Networth{
			Networth:     user.Cash + assetsWorth,
			Initiated_on: primitive.NewDateTimeFromTime(time.Now().UTC()),
			User_id:      user.ID,
		}
		db.GetNetworthCollection().InsertOne(context.TODO(), networth)
	}
}
