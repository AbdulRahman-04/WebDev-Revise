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
		ID:        primitive.NewObjectID(),
		UserId:    userId,
		FuncName:  funcName,
		FuncType:  funcType,
		FuncDesc:  funcDesc,
		IsPublic:  isPublic,
		ImageUrl: imageUrl,
		Status:    status,
		Location:  location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Goroutine for uploading file
	go func() {
		defer wg.Done()
		imageUrl, uploadErr = utils.FileUpload(c)
		if uploadErr != nil {
			imageUrl = ""
		}
	}()

	// Goroutine for DB insert
	go func() {
		defer wg.Done()
		_, insertErr = functionCollection.InsertOne(ctx, newFunction)
	}()

	wg.Wait()

	if insertErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	// Update image URL after insert (safe update)
	if uploadErr == nil && imageUrl != "" {
		update := bson.M{"$set": bson.M{"imageUrl": imageUrl, "updatedAt": time.Now()}}
		_, _ = functionCollection.UpdateByID(ctx, newFunction.ID, update)
		newFunction.ImageUrl = imageUrl
	}

	// Async Redis invalidation using background context
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			pattern := fmt.Sprintf("function_list:%s:*", userId.Hex())
			iter := utils.RedisClient.Scan(rctx, 0, pattern, 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}()

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

	cacheKey := fmt.Sprintf("function_list:%s:page:%d:limit:%d", userId.Hex(), page, limit)

	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var payload struct {
				Functions []models.Function `json:"functions"`
				Total     int64             `json:"total"`
			}
			if jerr := json.Unmarshal([]byte(cached), &payload); jerr == nil {
				c.JSON(200, gin.H{
					"msg":       "All Functions are here✨",
					"functions": payload.Functions,
					"page":      page,
					"limit":     limit,
					"total":     payload.Total,
					"source":    "redis",
				})
				return
			}
		}
	}

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
		cursor, err := functionCollection.Find(ctx, bson.M{"userId": userId}, opts)
		if err == nil {
			defer cursor.Close(ctx)
			findErr = cursor.All(ctx, &allFunctions)
		} else {
			findErr = err
		}
	}()

	wg.Wait()

	if countErr != nil || findErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	// Cache result async
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			payload := struct {
				Functions []models.Function `json:"functions"`
				Total     int64             `json:"total"`
			}{Functions: allFunctions, Total: total}
			if b, err := json.Marshal(payload); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 60*time.Second).Err()
			}
		}
	}()

	c.JSON(200, gin.H{
		"msg":       "All Functions are here✨",
		"functions": allFunctions,
		"page":      page,
		"limit":     limit,
		"total":     total,
		"source":    "db",
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

	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			if jerr := json.Unmarshal([]byte(cached), &oneFunc); jerr == nil {
				c.JSON(200, gin.H{"msg": "Function from Redis✅", "function": oneFunc, "source": "redis"})
				return
			}
		}
	}

	if err := functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneFunc); err != nil {
		c.JSON(404, gin.H{"msg": "No function found❌"})
		return
	}

	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			if b, err := json.Marshal(oneFunc); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 60*time.Second).Err()
			}
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
		c.JSON(404, gin.H{"msg": "No function found"})
		return
	}

	// fields
	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	// Channels for concurrency
	imageCh := make(chan string, 1)
	errCh := make(chan error, 2)
	done := make(chan struct{})

	// 🔹 1. Upload Image (if exists)
	go func() {
		defer close(imageCh)
		file, _ := c.FormFile("file")
		if file == nil {
			imageCh <- oldFunc.ImageUrl // keep old if not uploaded
			return
		}
		path, err := utils.FileUpload(c)
		if err != nil {
			errCh <- err
			imageCh <- oldFunc.ImageUrl
			return
		}
		imageCh <- path
	}()

	// 🔹 2. Update DB concurrently after image upload finishes
	go func() {
		imageUrl := <-imageCh // wait for upload result

		update := bson.M{"$set": bson.M{
			"funcname":  funcName,
			"functype":  funcType,
			"funcdes":   funcDesc,
			"ispublic":  isPublic,
			"status":    status,
			"location":  location,
			"imageUrl":  imageUrl,
			"updatedAt": time.Now(),
		}}

		setMap := update["$set"].(bson.M)
		for k, v := range setMap {
			if s, ok := v.(string); ok && s == "" {
				delete(setMap, k)
			}
		}

		if _, err := functionCollection.UpdateByID(ctx, oldFunc.ID, update); err != nil {
			errCh <- err
			close(done)
			return
		}

		// merge new values into oldFunc for response
		for k, v := range setMap {
			switch k {
			case "funcname":
				oldFunc.FuncName = v.(string)
			case "functype":
				oldFunc.FuncType = v.(string)
			case "funcdes":
				oldFunc.FuncDesc = v.(string)
			case "ispublic":
				oldFunc.IsPublic = v.(string)
			case "status":
				oldFunc.Status = v.(string)
			case "location":
				oldFunc.Location = v.(string)
			case "imageUrl":
				oldFunc.ImageUrl = v.(string)
			}
		}
		oldFunc.UpdatedAt = time.Now()

		// ✅ 3. Clear Redis Cache (async)
		go func() {
			rctx := context.Background()
			if utils.RedisClient != nil {
				cacheKey := fmt.Sprintf("function:%s:%s", userId.Hex(), oldFunc.ID.Hex())
				_ = utils.RedisClient.Del(rctx, cacheKey).Err()

				iter := utils.RedisClient.Scan(rctx, 0, fmt.Sprintf("function_list:%s:*", userId.Hex()), 0).Iterator()
				for iter.Next(rctx) {
					_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
				}
			}
		}()

		close(done)
	}()

	select {
	case <-done:
		c.JSON(200, gin.H{
			"msg":              "Function Updated Successfully ✅",
			"updatedFunction":  oldFunc,
			"response_time_ms": "Fast concurrent update",
		})
	case e := <-errCh:
		c.JSON(500, gin.H{"msg": "Error during update", "err": e.Error()})
	case <-ctx.Done():
		c.JSON(504, gin.H{"msg": "Request timed out"})
	}
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

	res, err := functionCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(404, gin.H{"msg": "No function found to delete"})
		return
	}

	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			cacheKey := fmt.Sprintf("function:%s:%s", userId.Hex(), mongoId.Hex())
			_ = utils.RedisClient.Del(rctx, cacheKey).Err()
			pattern := fmt.Sprintf("function_list:%s:*", userId.Hex())
			iter := utils.RedisClient.Scan(rctx, 0, pattern, 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
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

	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			patterns := []string{
				fmt.Sprintf("function_list:%s:*", userId.Hex()),
				fmt.Sprintf("function:%s:*", userId.Hex()),
			}
			for _, pattern := range patterns {
				iter := utils.RedisClient.Scan(rctx, 0, pattern, 0).Iterator()
				for iter.Next(rctx) {
					_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
				}
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "All Functions Deleted✅"})
}
