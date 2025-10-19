package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/AbdulRahman-04/07EvenetManagement-API/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func DbConnect() error {
	// ctx 
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connection 
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.DBURL))
    if err != nil {
		fmt.Println("Error connecting db")
	}

	// ping 
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("test db failed")
	}

	MongoClient = client
	fmt.Println("MongoDB Connected Successfull!`✅")
	return  nil
}