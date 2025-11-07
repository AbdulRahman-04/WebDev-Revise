package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JoinRequest model
type JoinRequest struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventID     primitive.ObjectID `bson:"eventId,omitempty" json:"eventId"`
	FunctionID  *primitive.ObjectID `bson:"functionId,omitempty" json:"functionId"`
	RequesterID primitive.ObjectID `bson:"requesterId,omitempty" json:"requesterId"`
	OwnerID     primitive.ObjectID `bson:"ownerId,omitempty" json:"ownerId"`
	Status      string             `bson:"status" json:"status"` // pending | accepted | rejected
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
