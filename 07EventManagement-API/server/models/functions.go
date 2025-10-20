package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Function struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id"`
	UserId primitive.ObjectID `bson:"user_id" json:"user_id"`

	FuncName string `bson:"funcName" json:"funcName" binding:"required,min=10,max=50"`
	FuncDesc string `bson:"funcDesc" json:"funcDesc" binding:"required,min=10,max=150"`
	FuncType string `bson:"funcType" json:"funcType" binding:"required,oneof=shadi valima wedding reception sanchak mehendi birthdayParty getTogether"`
	ImageURL string `bson:"imageURL" json:"imageURL" binding:"required"`
	IsPublic bool `bson:"isPublic" json:"isPublic" binding:"required"`
	Status string `bson:"status" json:"status" binding:"required,oneof=active closed comingSoon"`
	Location string `bson:"location" json:"location" binding:"required,min=15,max=100"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}