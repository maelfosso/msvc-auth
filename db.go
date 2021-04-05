package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context
var client *mongo.Client
var database *mongo.Database

func OpenDB() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}

	database = client.Database(os.Getenv("MONGODB_DBNAME"))
}

func Database() *mongo.Database {
	if database == nil {
		OpenDB()
	}
	return database
}

func CloseDB() {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
