package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id"`
	UserId primitive.ObjectID `bson:"user_id" json:"user_id"`

	EventName string `bson:"eventName" json:"eventName" binding:"required,min=10,max=50"`
	EventDesc string `bson:"eventDesc" json:"eventDesc" binding:"required,min=10,max=150"`
	ImageUrl string `bson:"ImageUrl" json:"ImageURL" binding:"rqeuired"`
    EventType string `bson:"eventType" json:"eventType" binding:"required"`
	IsPublic bool `bson:"isPublic" json:"isPublic" bidning:"required"`
	Status string `bson:"status" json:"status" binding:"required,oneof=active comingSoon closed"`
	Location string `bson:"location" json:"location" binding:"required,min=15,max=100"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
