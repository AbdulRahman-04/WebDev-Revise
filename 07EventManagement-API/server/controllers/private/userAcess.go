package private

import (
	"context"
	"sync"
	// "strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func UserAccessCollect() {
	userCollection = utils.MongoClient.Database("Event_Booking").Collection("user")
}

// getone user api
func GetOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid id format"})
		return
	}

	var oneUser models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneUser); err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Your Profile✨", "OneUser": oneUser})
}



// edit user api
func EditUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid id format"})
		return
	}

	tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	if tokenUserId.Hex() != mongoId.Hex() {
		c.JSON(403, gin.H{"msg": "Unauthorized: You can't touch other user's data❌"})
		return
	}

	var input struct {
		UserName string `json:"name" form:"name"`
		Language string `json:"language" form:"language"`
		Location string `json:"location" form:"location"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "invalid request"})
		return
	}

	var valid bool
	var updateErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		valid = input.UserName != "" && input.Language != "" && input.Location != ""
	}()

	go func() {
		defer wg.Done()
		update := bson.M{"$set": bson.M{
			"username":   input.UserName,
			"language":   input.Language,
			"location":   input.Location,
			"updated_at": time.Now(),
		}}
		_, updateErr = userCollection.UpdateByID(ctx, mongoId, update)
	}()

	wg.Wait()

	if !valid {
		c.JSON(400, gin.H{"msg": "Invalid Request, Please add some values to edit ur profile⚠️"})
		return
	}
	if updateErr != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Your Profile Updated Successfully!✨"})
}

// delete one api
func DeleteOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid id format"})
		return
	}

	tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	if tokenUserId.Hex() != mongoId.Hex() {
		c.JSON(403, gin.H{"msg": "Unauthorized: You can't touch other user's data❌"})
		return
	}

	if _, err := userCollection.DeleteOne(ctx, bson.M{"_id": mongoId}); err != nil {
		c.JSON(400, gin.H{"msg": "couldn't delete user, no id found!⚠️"})
		return
	}

	c.JSON(200, gin.H{"msg": "Your Profile Deleted💔"})
}


// User Logout API
func UserLogout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type LogoutInput struct {
		RefreshToken string `json:"refreshToken"`
	}

	var input LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	// Find user by refresh token
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"msg": "Invalid refresh token"})
		return
	}

	// Invalidate refresh token
	_, err = userCollection.UpdateByID(ctx, user.ID, bson.M{
		"$set": bson.M{
			"refreshToken":  "",
			"refreshExpiry": time.Time{},
			"updated_at":    time.Now(),
		},
	})

	if err != nil {
		c.JSON(500, gin.H{"msg": "Could not logout, try again"})
		return
	}

	c.JSON(200, gin.H{
		"msg": "User logged out successfully ✅",
	})
}
