package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := MongoClient.Database("Event_Booking")

	// =======================
	// USER COLLECTION INDEXES
	// =======================
	userIdx := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("email_unique"),
		},
		{
			Keys:    bson.D{{Key: "refreshToken", Value: 1}},
			Options: options.Index().SetName("refreshToken_idx"),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("user_createdAt_idx"),
		},
	}
	_, err := db.Collection("user").Indexes().CreateMany(ctx, userIdx)
	if err != nil {
		log.Println("❌ USER indexes error:", err)
	}

	// =======================
	// ADMIN COLLECTION INDEXES
	// =======================
	adminIdx := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("admin_email_unique"),
		},
		{
			Keys:    bson.D{{Key: "refreshToken", Value: 1}},
			Options: options.Index().SetName("admin_refreshToken_idx"),
		},
	}
	_, err = db.Collection("admin").Indexes().CreateMany(ctx, adminIdx)
	if err != nil {
		log.Println("❌ ADMIN indexes error:", err)
	}

	// =======================
	// EVENTS COLLECTION INDEXES
	// =======================
	eventIdx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}},
			Options: options.Index().SetName("userId_idx"),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("event_createdAt_idx"),
		},
	}
	_, err = db.Collection("events").Indexes().CreateMany(ctx, eventIdx)
	if err != nil {
		log.Println("❌ EVENTS indexes error:", err)
	}

	// =======================
	// FUNCTIONS COLLECTION INDEXES
	// =======================
	funcIdx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("function_createdAt_idx"),
		},
	}
	_, err = db.Collection("functions").Indexes().CreateMany(ctx, funcIdx)
	if err != nil {
		log.Println("❌ FUNCTIONS indexes error:", err)
	}

	// =======================
	// JOIN_REQUESTS INDEXES
	// =======================
	joinIdx := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "eventId", Value: 1}},
			Options: options.Index().SetName("eventId_idx"),
		},
		{
			Keys:    bson.D{{Key: "ownerId", Value: 1}},
			Options: options.Index().SetName("ownerId_idx"),
		},
		{
			Keys:    bson.D{{Key: "requesterId", Value: 1}},
			Options: options.Index().SetName("requesterId_idx"),
		},
		{
			Keys: bson.D{
				{Key: "eventId", Value: 1},
				{Key: "requesterId", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_event_requester"),
		},
	}
	_, err = db.Collection("join_requests").Indexes().CreateMany(ctx, joinIdx)
	if err != nil {
		log.Println("❌ JOIN_REQUEST indexes error:", err)
	}

	log.Println("✅ All MongoDB indexes created successfully!")
}
