package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var functionCollection *mongo.Collection

func FunctionCollect() {
	functionCollection = utils.MongoClient.Database("Event_Booking").Collection("functions")
}

func CreateFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	var imageUrl string
	var uploadErr, insertErr error

	newFunction := models.Function{
		ID:         primitive.NewObjectID(),
		UserId:     userId,
		FuncName:   funcName,
		FuncType:   funcType,
		FuncDesc:   funcDesc,
		IsPublic:   isPublic,
		Status:     status,
		Location:   location,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		imageUrl, uploadErr = utils.FileUpload(c)
		if uploadErr != nil {
			imageUrl = ""
		}
	}()

	go func() {
		defer wg.Done()
		_, insertErr = functionCollection.InsertOne(ctx, newFunction)
	}()

	wg.Wait()
	newFunction.ImageUrl = imageUrl

	if insertErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	c.JSON(200, gin.H{"msg": "New Function Created✨", "functionDetails": newFunction})
}

func GetAllFunctions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	var total int64
	var allFunctions []models.Function
	var countErr, findErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		total, countErr = functionCollection.CountDocuments(ctx, bson.M{"userId": userId})
	}()

	go func() {
		defer wg.Done()
		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "createdAt", Value: -1}})
		cursor, findErr := functionCollection.Find(ctx, bson.M{"userId": userId}, opts)
		if findErr == nil {
			defer cursor.Close(ctx)
			findErr = cursor.All(ctx, &allFunctions)
		}
	}()

	wg.Wait()
	if countErr != nil || findErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	c.JSON(200, gin.H{
		"msg":       "All Functions are here✨",
		"functions": allFunctions,
		"page":      page,
		"limit":     limit,
		"total":     total,
	})
}
func GetOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	cacheKey := fmt.Sprintf("function:%s:%s", userId.Hex(), mongoId.Hex())
	var oneFunc models.Function
	var redisErr, dbErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if utils.RedisClient != nil {
			if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
				redisErr = json.Unmarshal([]byte(cached), &oneFunc)
			}
		}
	}()

	go func() {
		defer wg.Done()
		dbErr = functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneFunc)
	}()

	wg.Wait()

	if redisErr == nil {
		c.JSON(200, gin.H{"msg": "Function from Redis✅", "function": oneFunc, "source": "redis"})
		return
	}
	if dbErr != nil {
		c.JSON(404, gin.H{"msg": "No function found❌"})
		return
	}

	go func() {
		if utils.RedisClient != nil {
			data, _ := json.Marshal(oneFunc)
			_ = utils.RedisClient.Set(ctx, cacheKey, data, 60*time.Second).Err()
		}
	}()

	c.JSON(200, gin.H{"msg": "Function from DB✅", "function": oneFunc, "source": "db"})
}

func EditFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	var oldFunc models.Function
	if err := functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oldFunc); err != nil {
		c.JSON(400, gin.H{"msg": "No function found to update"})
		return
	}

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	var imageUrl string
	var uploadErr, updateErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		imageUrl, uploadErr = utils.FileUpload(c)
		if uploadErr != nil {
			imageUrl = ""
		}
	}()

	go func() {
		defer wg.Done()
		update := bson.M{"$set": bson.M{
			"funcname":   funcName,
			"functype":   funcType,
			"funcdes":    funcDesc,
			"ispublic":   isPublic,
			"status":     status,
			"location":   location,
			"imageUrl":   imageUrl,
			"updated_at": time.Now(),
		}}
		_, updateErr = functionCollection.UpdateByID(ctx, oldFunc.ID, update)
	}()

	wg.Wait()
	if updateErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Function Updated Successfully!✅", "updatedFunction": oldFunc})
}

func DeleteOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	// Delete from DB
	_, err = functionCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "No function found to delete"})
		return
	}

	// Purge Redis cache (optional)
	go func() {
		cacheKey := fmt.Sprintf("function:%s:%s", userId.Hex(), mongoId.Hex())
		if utils.RedisClient != nil {
			_ = utils.RedisClient.Del(ctx, cacheKey).Err()
		}
	}()

	c.JSON(200, gin.H{"msg": "One Function Deleted✅"})
}

func DeleteAllFunctions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	_, err := functionCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	// Optional: purge Redis keys for this user
	go func() {
		if utils.RedisClient != nil {
			pattern := fmt.Sprintf("function:%s:*", userId.Hex())
			iter := utils.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
			for iter.Next(ctx) {
				_ = utils.RedisClient.Del(ctx, iter.Val()).Err()
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "All Functions Deleted✅"})
}