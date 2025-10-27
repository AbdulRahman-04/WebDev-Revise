package private

import (
	"context"
	"encoding/json"
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

// ========================
// 🔹 Global Collections
// ========================
var (
	UserCollection     *mongo.Collection
	EventCollection    *mongo.Collection
	FunctionCollection *mongo.Collection
	adminCollection    *mongo.Collection
)

// Connect Admin Collections
func AdminAccessCollect() {
	db := utils.MongoClient.Database("Event_Booking")
	adminCollection = db.Collection("admin")
	UserCollection = db.Collection("user")
	EventCollection = db.Collection("event")
	FunctionCollection = db.Collection("function")
}

// ========================================
// 🧑‍💼 GET ALL USERS (Admin) - userAccess style
// ========================================
func GetAllUsersAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 🔹 Direct pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit
	cacheKey := "admin:users:page=" + strconv.Itoa(page) + ":limit=" + strconv.Itoa(limit)

	// ✅ Redis cache first
	if cachedValue, err := utils.RedisGet(cacheKey); err == nil && cachedValue != "" {
		var users []models.User
		_ = json.Unmarshal([]byte(cachedValue), &users)

		total, _ := UserCollection.CountDocuments(ctx, bson.M{})
		totalPages := (total + int64(limit) - 1) / int64(limit)

		c.JSON(200, gin.H{
			"msg":       "All Users (from Redis)✨",
			"users":     users,
			"page":      page,
			"limit":     limit,
			"total":     total,
			"totalPage": totalPages,
			"hasNext":   page < int(totalPages),
			"hasPrev":   page > 1,
			"source":    "cache",
		})
		return
	}

	// ✅ Mongo fetch fallback
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := UserCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	_ = cursor.All(ctx, &users)

	total, _ := UserCollection.CountDocuments(ctx, bson.M{})
	totalPages := (total + int64(limit) - 1) / int64(limit)

	data, _ := json.Marshal(users)
	_ = utils.RedisSet(cacheKey, string(data))

	c.JSON(200, gin.H{
		"msg":       "All Users✨",
		"users":     users,
		"page":      page,
		"limit":     limit,
		"total":     total,
		"totalPage": totalPages,
		"hasNext":   page < int(totalPages),
		"hasPrev":   page > 1,
		"source":    "db",
	})
}

// ========================================
// 🧑‍💼 GET ONE USER - userAccess style
// ========================================
func GetOneUserAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := "admin:user:" + id
	if cachedValue, err := utils.RedisGet(cacheKey); err == nil && cachedValue != "" {
		var user models.User
		_ = json.Unmarshal([]byte(cachedValue), &user)
		c.JSON(200, gin.H{"msg": "User (from Redis)✨", "user": user, "source": "cache"})
		return
	}

	var user models.User
	if err := UserCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&user); err != nil {
		c.JSON(404, gin.H{"msg": "User not found❌"})
		return
	}

	data, _ := json.Marshal(user)
	_ = utils.RedisSet(cacheKey, string(data))

	c.JSON(200, gin.H{"msg": "User✨", "user": user, "source": "db"})
}

// ========================================
// 🎉 GET ALL EVENTS - userAccess style
// ========================================
func GetAllEventsAdmin(c *gin.Context) {
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
	cacheKey := "admin:events:page=" + strconv.Itoa(page) + ":limit=" + strconv.Itoa(limit)

	if cachedValue, err := utils.RedisGet(cacheKey); err == nil && cachedValue != "" {
		var events []models.Event
		_ = json.Unmarshal([]byte(cachedValue), &events)

		total, _ := EventCollection.CountDocuments(ctx, bson.M{})
		totalPages := (total + int64(limit) - 1) / int64(limit)

		c.JSON(200, gin.H{
			"msg":       "All Events (from Redis)✨",
			"events":    events,
			"page":      page,
			"limit":     limit,
			"total":     total,
			"totalPage": totalPages,
			"hasNext":   page < int(totalPages),
			"hasPrev":   page > 1,
			"source":    "cache",
		})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := EventCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error"})
		return
	}
	defer cursor.Close(ctx)

	var events []models.Event
	_ = cursor.All(ctx, &events)

	total, _ := EventCollection.CountDocuments(ctx, bson.M{})
	totalPages := (total + int64(limit) - 1) / int64(limit)

	data, _ := json.Marshal(events)
	_ = utils.RedisSet(cacheKey, string(data))

	c.JSON(200, gin.H{
		"msg":       "All Events✨",
		"events":    events,
		"page":      page,
		"limit":     limit,
		"total":     total,
		"totalPage": totalPages,
		"hasNext":   page < int(totalPages),
		"hasPrev":   page > 1,
		"source":    "db",
	})
}

// ========================================
// 🎉 GET ONE EVENT - userAccess style
// ========================================
func GetOneEventAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := "admin:event:" + id
	if cachedValue, err := utils.RedisGet(cacheKey); err == nil && cachedValue != "" {
		var event models.Event
		_ = json.Unmarshal([]byte(cachedValue), &event)
		c.JSON(200, gin.H{"msg": "Event (from Redis)✨", "event": event, "source": "cache"})
		return
	}

	var event models.Event
	if err := EventCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&event); err != nil {
		c.JSON(404, gin.H{"msg": "Event not found❌"})
		return
	}

	data, _ := json.Marshal(event)
	_ = utils.RedisSet(cacheKey, string(data))

	c.JSON(200, gin.H{"msg": "Event✨", "event": event, "source": "db"})
}

// ========================================
// 🎭 GET ALL FUNCTIONS - userAccess style
// ========================================
func GetAllFunctionsAdmin(c *gin.Context) {
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
	cacheKey := "admin:functions:page=" + strconv.Itoa(page) + ":limit=" + strconv.Itoa(limit)

	if cachedValue, err := utils.RedisGet(cacheKey); err == nil && cachedValue != "" {
		var funcs []models.Function
		_ = json.Unmarshal([]byte(cachedValue), &funcs)

		total, _ := FunctionCollection.CountDocuments(ctx, bson.M{})
		totalPages := (total + int64(limit) - 1) / int64(limit)

		c.JSON(200, gin.H{
			"msg":       "All Functions (from Redis)✨",
			"functions": funcs,
			"page":      page,
			"limit":     limit,
			"total":     total,
			"totalPage": totalPages,
			"hasNext":   page < int(totalPages),
			"hasPrev":   page > 1,
			"source":    "cache",
		})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := FunctionCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error"})
		return
	}
	defer cursor.Close(ctx)

	var funcs []models.Function
	_ = cursor.All(ctx, &funcs)

	total, _ := FunctionCollection.CountDocuments(ctx, bson.M{})
	totalPages := (total + int64(limit) - 1) / int64(limit)

	data, _ := json.Marshal(funcs)
	_ = utils.RedisSet(cacheKey, string(data))

	c.JSON(200, gin.H{
		"msg":       "All Functions✨",
		"functions": funcs,
		"page":      page,
		"limit":     limit,
		"total":     total,
		"totalPage": totalPages,
		"hasNext":   page < int(totalPages),
		"hasPrev":   page > 1,
		"source":    "db",
	})
}

// ========================================
// 🎭 GET ONE FUNCTION - userAccess style
// ========================================
func GetOneFunctionAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := "admin:function:" + id
	if cachedValue, err := utils.RedisGet(cacheKey); err == nil && cachedValue != "" {
		var f models.Function
		_ = json.Unmarshal([]byte(cachedValue), &f)
		c.JSON(200, gin.H{"msg": "Function (from Redis)✨", "function": f, "source": "cache"})
		return
	}

	var f models.Function
	if err := FunctionCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&f); err != nil {
		c.JSON(404, gin.H{"msg": "Function not found❌"})
		return
	}

	data, _ := json.Marshal(f)
	_ = utils.RedisSet(cacheKey, string(data))

	c.JSON(200, gin.H{"msg": "Function✨", "function": f, "source": "db"})
}

// ========================================
// 🔒 Admin Logout
// ========================================
func AdminLogout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		c.JSON(400, gin.H{"msg": "Invalid request❌"})
		return
	}

	var admin models.Admin
	if err := adminCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&admin); err != nil {
		c.JSON(401, gin.H{"msg": "Invalid refresh token❌"})
		return
	}

	_, err := adminCollection.UpdateByID(ctx, admin.ID, bson.M{
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

	_ = utils.RedisDel("admin:" + admin.ID.Hex())
	c.JSON(200, gin.H{"msg": "Admin logged out successfully ✅"})
}
