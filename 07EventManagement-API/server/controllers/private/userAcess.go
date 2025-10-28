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

// ====================================================
// ⚡ Get All Users (Fast + Redis + Pagination)
// ====================================================
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

	type CachedData struct {
		Users []models.User `json:"users"`
		Total int64         `json:"total"`
	}

	// ✅ Try Redis cache first
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var payload CachedData
			if err := json.Unmarshal([]byte(cached), &payload); err == nil {
				c.JSON(200, gin.H{
					"msg":    "All Users (from Redis Cache)✨",
					"users":  payload.Users,
					"total":  payload.Total,
					"page":   page,
					"limit":  limit,
					"source": "redis",
				})
				return
			}
		}
	}

	// ✅ Fetch from MongoDB
	var users []models.User
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdat", Value: -1}}).
		SetBatchSize(20)
	cursor, err := userCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "Database error❌"})
		return
	}
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"msg": "Cursor decode error❌"})
		return
	}
	total, _ := userCollection.CountDocuments(ctx, bson.M{})

	// ✅ Cache result asynchronously
	if utils.RedisClient != nil {
		go func(users []models.User, total int64) {
			rctx := context.Background()
			payload := CachedData{Users: users, Total: total}
			if b, err := json.Marshal(payload); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 12*time.Hour).Err()
			}
		}(users, total)
	}

	c.JSON(200, gin.H{
		"msg":    "All Users (from DB)✨",
		"users":  users,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"source": "db",
	})
}

// ====================================================
// ⚡ Get One User (Fast + Redis Cache)
// ====================================================
func GetOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid user ID❌"})
		return
	}
	cacheKey := fmt.Sprintf("user:%s", mongoId.Hex())

	// ✅ Try Redis first
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var user models.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				c.JSON(200, gin.H{"msg": "User from Redis✅", "user": user, "source": "redis"})
				return
			}
		}
	}

	// ✅ Fallback: MongoDB
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&user); err != nil {
		c.JSON(404, gin.H{"msg": "User not found❌"})
		return
	}

	// ✅ Cache async
	if utils.RedisClient != nil {
		go func(user models.User) {
			rctx := context.Background()
			if b, err := json.Marshal(user); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 12*time.Hour).Err()
			}
		}(user)
	}

	c.JSON(200, gin.H{"msg": "User from DB✅", "user": user, "source": "db"})
}

// ====================================================
// ⚡ Edit User (Invalidate Redis Async)
// ====================================================
func EditUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid user ID❌"})
		return
	}

	tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	if tokenUserId.Hex() != mongoId.Hex() {
		c.JSON(403, gin.H{"msg": "Unauthorized❌"})
		return
	}

	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid JSON❌"})
		return
	}

	updateFields := bson.M{}
	for k, v := range input {
		if v != "" {
			updateFields[k] = v
		}
	}
	if len(updateFields) == 0 {
		c.JSON(400, gin.H{"msg": "No valid fields to update❌"})
		return
	}

	updateFields["updatedat"] = time.Now()
	_, err = userCollection.UpdateByID(ctx, mongoId, bson.M{"$set": updateFields})
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB update error❌"})
		return
	}

	// ✅ Invalidate Redis (background)
	if utils.RedisClient != nil {
		go func(id primitive.ObjectID) {
			rctx := context.Background()
			_ = utils.RedisClient.Del(rctx, fmt.Sprintf("user:%s", id.Hex())).Err()
			iter := utils.RedisClient.Scan(rctx, 0, "user_list:*", 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}(mongoId)
	}

	c.JSON(200, gin.H{"msg": "Profile Updated Successfully✅", "updated": updateFields})
}

// ====================================================
// ⚡ Delete User (Invalidate Redis Async)
// ====================================================
func DeleteOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid user ID❌"})
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

	// ✅ Invalidate Redis (background)
	if utils.RedisClient != nil {
		go func(id primitive.ObjectID) {
			rctx := context.Background()
			_ = utils.RedisClient.Del(rctx, fmt.Sprintf("user:%s", id.Hex())).Err()
			iter := utils.RedisClient.Scan(rctx, 0, "user_list:*", 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}(mongoId)
	}

	c.JSON(200, gin.H{"msg": "User Deleted Successfully💔"})
}

// ====================================================
// ⚡ User Logout (Clear Token + Redis)
// ====================================================
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

	// ✅ Remove Redis cache (background)
	if utils.RedisClient != nil {
		go func(id primitive.ObjectID) {
			rctx := context.Background()
			_ = utils.RedisClient.Del(rctx, "user:"+id.Hex()).Err()
		}(user.ID)
	}

	c.JSON(200, gin.H{"msg": "User logged out successfully ✅"})
}
