package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func UserAccessCollect() {
	userCollection = utils.MongoClient.Database("Event_Booking").Collection("user")
}

// ========================================
// 🧩 Get All Users (Redis + Pagination same as funcs/events)
// ========================================
func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	cacheKey := fmt.Sprintf("user_list:page:%d:limit:%d", page, limit)

	// ✅ Try Redis cache first
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var payload struct {
				Users []models.User `json:"users"`
				Total int64         `json:"total"`
			}
			if jerr := json.Unmarshal([]byte(cached), &payload); jerr == nil {
				c.JSON(200, gin.H{
					"msg":       "All Users (from Redis Cache)✨",
					"users":     payload.Users,
					"page":      page,
					"limit":     limit,
					"total":     payload.Total,
					"source":    "redis",
				})
				return
			}
		}
	}

	// ✅ Fetch from MongoDB
	var total int64
	var users []models.User

	total, _ = userCollection.CountDocuments(ctx, bson.M{})
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := userCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "Database error"})
		return
	}
	defer cursor.Close(ctx)
	_ = cursor.All(ctx, &users)

	// ✅ Cache the result (24 hours expiry)
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			payload := struct {
				Users []models.User `json:"users"`
				Total int64         `json:"total"`
			}{Users: users, Total: total}
			if b, err := json.Marshal(payload); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 24*time.Hour).Err()
			}
		}
	}()

	c.JSON(200, gin.H{
		"msg":       "All Users Are Here✨",
		"users":     users,
		"page":      page,
		"limit":     limit,
		"total":     total,
		"source":    "db",
	})
}

// ========================================
// 🧩 Get One User (Redis cache same as funcs/events)
// ========================================
func GetOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := fmt.Sprintf("user:%s", mongoId.Hex())
	var oneUser models.User

	// ✅ Check Redis cache
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			if jerr := json.Unmarshal([]byte(cached), &oneUser); jerr == nil {
				c.JSON(200, gin.H{"msg": "User from Redis✅", "user": oneUser, "source": "redis"})
				return
			}
		}
	}

	// ✅ Fetch from DB
	if err := userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneUser); err != nil {
		c.JSON(404, gin.H{"msg": "No user found❌"})
		return
	}

	// ✅ Cache in Redis for 24h
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			if b, err := json.Marshal(oneUser); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 24*time.Hour).Err()
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "User from DB✅", "user": oneUser, "source": "db"})
}

func EditUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID format"})
		return
	}

	tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	if tokenUserId.Hex() != mongoId.Hex() {
		c.JSON(403, gin.H{"msg": "Unauthorized: You can't edit other user's data❌"})
		return
	}

	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	// 🟢 Get existing user
	var existingUser models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&existingUser); err != nil {
		c.JSON(404, gin.H{"msg": "User not found"})
		return
	}

	// 🟢 Prepare update map (ignore empty string)
	updateFields := bson.M{}
	if v, ok := input["username"]; ok && v != "" {
		updateFields["username"] = v
	}
	if v, ok := input["language"]; ok && v != "" {
		updateFields["language"] = v
	}
	if v, ok := input["location"]; ok && v != "" {
		updateFields["location"] = v
	}
	if v, ok := input["phone"]; ok && v != "" {
		updateFields["phone"] = v
	}

	if len(updateFields) == 0 {
		c.JSON(400, gin.H{"msg": "No valid fields to update"})
		return
	}

	updateFields["updatedat"] = time.Now()
	update := bson.M{"$set": updateFields}

	_, err = userCollection.UpdateByID(ctx, mongoId, update)
	if err != nil {
		c.JSON(500, gin.H{"msg": "Database update error"})
		return
	}

	// 🟢 Update Redis cache
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			cacheKey := fmt.Sprintf("user:%s", mongoId.Hex())
			_ = utils.RedisClient.Del(rctx, cacheKey).Err()

			iter := utils.RedisClient.Scan(rctx, 0, "user_list:*", 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "Profile Updated Successfully!✅", "updatedFields": updateFields})
}




// ========================================
// 🧩 Delete One User (invalidate cache)
// ========================================
func DeleteOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	if userId.Hex() != mongoId.Hex() {
		c.JSON(403, gin.H{"msg": "Unauthorized❌"})
		return
	}

	res, err := userCollection.DeleteOne(ctx, bson.M{"_id": mongoId})
	if err != nil || res.DeletedCount == 0 {
		c.JSON(404, gin.H{"msg": "User not found or already deleted⚠️"})
		return
	}

	// ✅ Invalidate Redis
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			_ = utils.RedisClient.Del(rctx, fmt.Sprintf("user:%s", mongoId.Hex())).Err()
			iter := utils.RedisClient.Scan(rctx, 0, "user_list:*", 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "User Deleted Successfully💔"})
}
// ========================================
// 🔒 User Logout
// ========================================
func UserLogout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		c.JSON(400, gin.H{"msg": "Invalid request❌"})
		return
	}

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&user); err != nil {
		c.JSON(401, gin.H{"msg": "Invalid refresh token❌"})
		return
	}

	// ✅ Remove token from DB
	_, err := userCollection.UpdateByID(ctx, user.ID, bson.M{
		"$set": bson.M{
			"refreshToken":  "",
			"refreshExpiry": time.Time{},
			"updated_at":    time.Now(),
		},
	})
	if err != nil {
		c.JSON(500, gin.H{"msg": "Logout failed❌"})
		return
	}

	// ✅ Delete Redis cache
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			_ = utils.RedisClient.Del(rctx, "user:"+user.ID.Hex()).Err()
		}
	}()

	c.JSON(200, gin.H{"msg": "User logged out successfully ✅"})
}
