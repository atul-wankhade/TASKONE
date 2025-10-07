package db

import (
	"TASKONE/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		log.Fatal("Failed to connect to mongo DB", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB not reachable", err)
	}

	MongoClient = client
	MongoDatabase = client.Database(config.AppConfig.MongoDB)
	log.Println("Connected to mongoDB")
}
