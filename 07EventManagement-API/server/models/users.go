package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Role string `bson:"role" json:"role"`

	UserName string `bson:"userName" json:"userName" binding:"required"`
	Email string `bson:"email" json:"email" binding:"required,email"`
	Password string `bson:"password" json:"password" binding:"required,min=6,max=12"`
	Phone string `bson:"phone" json:"phone" binding:"required,min=10,max-15"`
	Language string `bson:"language" json:"language" binding:"required,oneof=hindi english urdu arabic tamil"`
	UserVerified struct {
		Email bool `bson:"emailVerified" json:"emailVerified"`
	} `bson:"emailVerified" json:"emailVerified"`

	UserVerifyToken struct {
		Email string `bson:"userVerifyToken" json:"userVerifyToken"`
	} `bson:"userVerifyToken" json:"userVerifyToken"`

	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
	RefreshExpiry time.Time `bson:"refreshExpiry" json:"refreshExpiry"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

}