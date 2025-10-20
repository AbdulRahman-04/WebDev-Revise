package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Role string `bson:"role" json:"role"`

	AdminName string `bson:"amdin_name" json:"admin_name" binding:"required"`
	Email string `bson:"email" json:"email" binding:"required,email"`
	Password string `bson:"password" json:"password" binding:"required,min=6,max=12"`
	Phone string `bson:"phone" json:"phone" binding:"required,min=10,max=15"`
	Language string `bson:"language" json:"language" binding:"required,oneof=hindi english urdu arabic tamil"`

	AdminVerified struct {
		Email bool `bson:"emailVerified" json:"emailVerified"`
	} `bson:"adminVerified" json:"adminVerified"`

	AdminVerifyToken struct {
		Email string `bson:"emailVerifyToken" json:"emailVerifyToken"`
	} `bson:"adminVerifyToken" json:"adminVerifyToken"`

	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
	RefreshExpiry time.Time `bson:"refreshExpiry" json:"refreshExpiry"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}