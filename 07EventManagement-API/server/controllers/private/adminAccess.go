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

var (
	UserCollection     *mongo.Collection
	EventCollection    *mongo.Collection
	FunctionCollection *mongo.Collection
	AdminCollection    *mongo.Collection
)

func AdminAccessCollect() {
	db := utils.MongoClient.Database("Event_Booking")
	AdminCollection = db.Collection("admin")
	UserCollection = db.Collection("user")
	EventCollection = db.Collection("events")
	FunctionCollection = db.Collection("functions")
}

// ==========================
// 🔥 GET ALL USERS (ADMIN)
// ==========================
func GetAllUsersAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}

	skip := (page - 1) * limit
	cacheKey := fmt.Sprintf("admin:users:%d:%d", page, limit)

	// Redis Check
	if cached, _ := utils.RedisClient.Get(ctx, cacheKey).Result(); cached != "" {
		var users []models.User
		_ = json.Unmarshal([]byte(cached), &users)

		for i := range users {
			users[i].Password = ""
			users[i].RefreshToken = ""
		}

		c.JSON(200, gin.H{
			"msg":    "Users (cache)✨",
			"users":  users,
			"page":   page,
			"limit":  limit,
			"source": "cache",
		})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cursor, err := UserCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error❌"})
		return
	}

	var users []models.User
	_ = cursor.All(ctx, &users)

	for i := range users {
		users[i].Password = ""
		users[i].RefreshToken = ""
	}

	// Cache it
	go utils.RedisClient.Set(ctx, cacheKey, toJson(users), 6*time.Hour)

	c.JSON(200, gin.H{"msg": "Users (db)✨", "users": users, "source": "db"})
}

// ==========================
// 🔥 GET ONE USER
// ==========================
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

	if cached, _ := utils.RedisClient.Get(ctx, cacheKey).Result(); cached != "" {
		var user models.User
		_ = json.Unmarshal([]byte(cached), &user)

		user.Password = ""
		user.RefreshToken = ""

		c.JSON(200, gin.H{"msg": "User (cache)✨", "user": user})
		return
	}

	var user models.User
	if err := UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		c.JSON(404, gin.H{"msg": "User not found❌"})
		return
	}

	user.Password = ""
	user.RefreshToken = ""

	go utils.RedisClient.Set(ctx, cacheKey, toJson(user), 6*time.Hour)

	c.JSON(200, gin.H{"msg": "User✨", "user": user})
}

// ==========================
// 🔥 GET ALL EVENTS
// ==========================
func AdminGetAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}

	skip := (page - 1) * limit
	cacheKey := fmt.Sprintf("admin:events:%d:%d", page, limit)

	if cached, _ := utils.RedisClient.Get(ctx, cacheKey).Result(); cached != "" {
		var events []models.Event
		_ = json.Unmarshal([]byte(cached), &events)

		c.JSON(200, gin.H{"msg": "Events (cache)✨", "events": events})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})

	cursor, err := EventCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error❌"})
		return
	}

	var events []models.Event
	_ = cursor.All(ctx, &events)

	go utils.RedisClient.Set(ctx, cacheKey, toJson(events), 10*time.Minute)

	c.JSON(200, gin.H{"msg": "Events (db)✨", "events": events})
}

// ==========================
// 🔥 GET ONE EVENT
// ==========================
func GetOneEventAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	var event models.Event
	if err := EventCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event); err != nil {
		c.JSON(404, gin.H{"msg": "Event not found❌"})
		return
	}

	c.JSON(200, gin.H{"msg": "Event✨", "event": event})
}

// ==========================
// 🔥 GET ALL FUNCTIONS
// ==========================
func GetAllFunctionsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := FunctionCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error❌"})
		return
	}

	var funcs []models.Function
	_ = cursor.All(ctx, &funcs)

	c.JSON(200, gin.H{"msg": "Functions✨", "functions": funcs})
}

// ==========================
// 🔥 GET ONE FUNCTION
// ==========================
func GetOneFunctionAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid ID❌"})
		return
	}

	var f models.Function
	if err := FunctionCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&f); err != nil {
		c.JSON(404, gin.H{"msg": "Function not found❌"})
		return
	}

	c.JSON(200, gin.H{"msg": "Function✨", "function": f})
}

// ==========================
// 🔒 LOGOUT ADMIN
// ==========================
func AdminLogout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payload struct {
		RefreshToken string `json:"refreshToken"`
	}
	_ = c.ShouldBindJSON(&payload)

	AdminCollection.UpdateOne(ctx, bson.M{"refreshToken": payload.RefreshToken}, bson.M{
		"$set": bson.M{"refreshToken": "", "refreshExpiry": time.Time{}},
	})

	c.JSON(200, gin.H{"msg": "Admin logged out successfully ✅"})
}

// helper
func toJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
