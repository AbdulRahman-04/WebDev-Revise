package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func DBConnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := config.AppConfig.DBURI
	fmt.Println("🔍 MONGO_URI from config:", uri)

	if uri == "" {
		log.Fatal("❌ MONGO_URI missing from .env")
		return fmt.Errorf("MONGO_URI missing")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("❌ Couldn't connect to MongoDB: %v", err)
	}

	// nil safety check
	if client == nil {
		return fmt.Errorf("❌ Mongo client is nil – check your URI")
	}

	// Ping test
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("❌ Test request failed: %v", err)
	}

	MongoClient = client
	fmt.Println("✅ MongoDB Connected Successfully!")
	return nil
}
