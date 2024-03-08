package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectToDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	connection, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	client = connection
}

func GetUserCollection() *mongo.Collection {
	return client.Database("StockPaperTradingApp").Collection("User")
}

func GetHoldingsCollection() *mongo.Collection {
	return client.Database("StockPaperTradingApp").Collection("Holdings")
}

func GetActivityCollection() *mongo.Collection {
	return client.Database("StockPaperTradingApp").Collection("Activity")
}

func GetNetworthCollection() *mongo.Collection {
	return client.Database("StockPaperTradingApp").Collection("Networth")
}
