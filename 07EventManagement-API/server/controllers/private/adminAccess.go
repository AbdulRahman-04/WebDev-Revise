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

// ========================
// 🔹 Global Collections
// ========================
var (
	UserCollection      *mongo.Collection
	EventCollection     *mongo.Collection
	FunctionCollection  *mongo.Collection
	AdminCollection     *mongo.Collection
)

// ========================
// 🔹 Connect Admin Collections
// ========================
func AdminAccessCollect() {
	db := utils.MongoClient.Database("Event_Booking")
	AdminCollection = db.Collection("admin")
	UserCollection = db.Collection("user")
	EventCollection = db.Collection("events")
	FunctionCollection = db.Collection("functions")
}

// ========================================
// 🧑‍💼 GET ALL USERS (Admin)
// ========================================
func GetAllUsersAdmin(c *gin.Context) {
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

	cacheKey := "admin:users:page=" + strconv.Itoa(page) + ":limit=" + strconv.Itoa(limit)
	if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var users []models.User
		_ = json.Unmarshal([]byte(cached), &users)
		total, _ := UserCollection.CountDocuments(ctx, bson.M{})
		totalPages := (total + int64(limit) - 1) / int64(limit)
		c.JSON(200, gin.H{
			"msg": "All Users (from Redis)✨", "users": users,
			"page": page, "limit": limit, "total": total, "totalPage": totalPages,
			"hasNext": page < int(totalPages), "hasPrev": page > 1, "source": "cache",
		})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := UserCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error❌"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"msg": "Error decoding users❌"})
		return
	}

	if len(users) > 0 {
		data, _ := json.Marshal(users)
		_ = utils.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	}

	total, _ := UserCollection.CountDocuments(ctx, bson.M{})
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(200, gin.H{
		"msg": "All Users✨", "users": users,
		"page": page, "limit": limit, "total": total, "totalPage": totalPages,
		"hasNext": page < int(totalPages), "hasPrev": page > 1, "source": "db",
	})
}

// ========================================
// 🧑‍💼 GET ONE USER (Admin)
// ========================================
func GetOneUserAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := "admin:user:" + id
	if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var user models.User
		_ = json.Unmarshal([]byte(cached), &user)
		c.JSON(200, gin.H{"msg": "User (from Redis)✨", "user": user, "source": "cache"})
		return
	}

	var user models.User
	if err := UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		c.JSON(404, gin.H{"msg": "User not found❌"})
		return
	}

	data, _ := json.Marshal(user)
	_ = utils.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()

	c.JSON(200, gin.H{"msg": "User✨", "user": user, "source": "db"})
}

// ========================================

func AdminGetAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}

	skip := (page - 1) * limit
	cacheKey := fmt.Sprintf("admin_events:page=%d:limit=%d", page, limit)

	// Try Redis cache first
	if utils.RedisClient != nil {
		cached, err := utils.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var cachedData struct {
				Events    []models.Event `json:"events"`
				Page      int            `json:"page"`
				Limit     int            `json:"limit"`
				Total     int64          `json:"total"`
				TotalPage int            `json:"totalPage"`
				HasNext   bool           `json:"hasNext"`
				HasPrev   bool           `json:"hasPrev"`
			}

			if err := json.Unmarshal([]byte(cached), &cachedData); err == nil {
				c.JSON(200, gin.H{
					"msg":       "All Events (from Redis)✨",
					"events":    cachedData.Events,
					"page":      cachedData.Page,
					"limit":     cachedData.Limit,
					"total":     cachedData.Total,
					"totalPage": cachedData.TotalPage,
					"hasNext":   cachedData.HasNext,
					"hasPrev":   cachedData.HasPrev,
					"source":    "cache",
				})
				return
			}
		}
	}

	// MongoDB query
	count, err := eventsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"msg": "Count error"})
		return
	}

	cur, err := eventsCollection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB Find error"})
		return
	}
	defer cur.Close(ctx)

	var events []models.Event
	if err := cur.All(ctx, &events); err != nil {
		c.JSON(500, gin.H{"msg": "Cursor error"})
		return
	}

	totalPages := int((count + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	responseData := gin.H{
		"msg":       "All Events (from MongoDB)✨",
		"events":    events,
		"page":      page,
		"limit":     limit,
		"total":     count,
		"totalPage": totalPages,
		"hasNext":   hasNext,
		"hasPrev":   hasPrev,
		"source":    "db",
	}

	// Cache store
	if utils.RedisClient != nil && len(events) > 0 {
		jsonData, _ := json.Marshal(responseData)
		_ = utils.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err()
	}

	c.JSON(200, responseData)
}



// ========================================
// 🎉 GET ONE EVENT (Admin)
// ========================================
func GetOneEventAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := "admin:event:" + id
	if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var event models.Event
		_ = json.Unmarshal([]byte(cached), &event)
		c.JSON(200, gin.H{"msg": "Event (from Redis)✨", "event": event, "source": "cache"})
		return
	}

	var event models.Event
	if err := EventCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event); err != nil {
		c.JSON(404, gin.H{"msg": "Event not found❌"})
		return
	}

	data, _ := json.Marshal(event)
	_ = utils.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()

	c.JSON(200, gin.H{"msg": "Event✨", "event": event, "source": "db"})
}

// ========================================
// 🎭 GET ALL FUNCTIONS (Admin)
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

	if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var funcs []models.Function
		_ = json.Unmarshal([]byte(cached), &funcs)
		total, _ := FunctionCollection.CountDocuments(ctx, bson.M{})
		totalPages := (total + int64(limit) - 1) / int64(limit)
		c.JSON(200, gin.H{
			"msg": "All Functions (from Redis)✨", "functions": funcs,
			"page": page, "limit": limit, "total": total, "totalPage": totalPages,
			"hasNext": page < int(totalPages), "hasPrev": page > 1, "source": "cache",
		})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := FunctionCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error❌"})
		return
	}
	defer cursor.Close(ctx)

	var funcs []models.Function
	if err := cursor.All(ctx, &funcs); err != nil {
		c.JSON(500, gin.H{"msg": "Error decoding functions❌"})
		return
	}

	if len(funcs) > 0 {
		data, _ := json.Marshal(funcs)
		_ = utils.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	}

	total, _ := FunctionCollection.CountDocuments(ctx, bson.M{})
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(200, gin.H{
		"msg": "All Functions✨", "functions": funcs,
		"page": page, "limit": limit, "total": total, "totalPage": totalPages,
		"hasNext": page < int(totalPages), "hasPrev": page > 1, "source": "db",
	})
}

// ========================================
// 🎭 GET ONE FUNCTION (Admin)
// ========================================
func GetOneFunctionAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	cacheKey := "admin:function:" + id
	if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var f models.Function
		_ = json.Unmarshal([]byte(cached), &f)
		c.JSON(200, gin.H{"msg": "Function (from Redis)✨", "function": f, "source": "cache"})
		return
	}

	var f models.Function
	if err := FunctionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&f); err != nil {
		c.JSON(404, gin.H{"msg": "Function not found❌"})
		return
	}

	data, _ := json.Marshal(f)
	_ = utils.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute).Err()

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
	if err := AdminCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&admin); err != nil {
		c.JSON(401, gin.H{"msg": "Invalid refresh token❌"})
		return
	}

	_, err := AdminCollection.UpdateByID(ctx, admin.ID, bson.M{
		"$set": bson.M{
			"refreshToken":  "",
			"refreshExpiry": time.Time{},
			"updatedAt":     time.Now(),
		},
	})
	if err != nil {
		c.JSON(500, gin.H{"msg": "Logout failed❌"})
		return
	}

	_ = utils.RedisClient.Del(ctx, "admin:"+admin.ID.Hex()).Err()
	c.JSON(200, gin.H{"msg": "Admin logged out successfully ✅"})
}
