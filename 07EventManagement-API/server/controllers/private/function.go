package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// -------------------- HELPER: Privacy apply on single function --------------------
func applyFunctionPrivacy(fn models.Function, userId primitive.ObjectID) models.Function {
	val := strings.ToLower(fn.IsPublic)
	isPublic := val == "public" || val == "true"
	isOwner := fn.UserId.Hex() == userId.Hex()

	// Public or owner → full data allowed
	if isPublic || isOwner {
		return fn
	}

	// Private & NOT owner → mask sensitive fields
	fn.Location = ""
	fn.ImageUrl = ""
	fn.Status = "Private"
	fn.FuncDesc = "This is a private function. Request to join to see full details."

	return fn
}

// -------------------- CREATE FUNCTION --------------------
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

	// Step 1: Upload file (if any)
	imageUrl := ""
	if file, _ := c.FormFile("file"); file != nil {
		if path, err := utils.FileUpload(c); err == nil {
			imageUrl = path
		}
	}

	// Step 2: Prepare new function
	newFunction := models.Function{
		ID:        primitive.NewObjectID(),
		UserId:    userId,
		FuncName:  funcName,
		FuncType:  funcType,
		FuncDesc:  funcDesc,
		IsPublic:  isPublic,
		ImageUrl:  imageUrl,
		Status:    status,
		Location:  location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Step 3: Insert DB
	if _, err := functionCollection.InsertOne(ctx, newFunction); err != nil {
		c.JSON(500, gin.H{"msg": "DB insert failed", "error": err.Error()})
		return
	}

	// Step 4: Clear cache async
	go func(uid string) {
		rctx := context.Background()
		if utils.RedisClient != nil {
			pattern := fmt.Sprintf("function_list:%s:*", uid)
			iter := utils.RedisClient.Scan(rctx, 0, pattern, 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}(userId.Hex())

	c.JSON(200, gin.H{
		"msg":             "New Function Created ✅",
		"functionDetails": newFunction,
	})
}

// -------------------- GET ALL FUNCTIONS (with privacy) --------------------
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

	// cache per-user (kyunki privacy user-specific hai)
	cacheKey := fmt.Sprintf("function_list:%s:page:%d:limit:%d", userId.Hex(), page, limit)

	// Try Redis
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

	// Count ALL functions (not just user ke)
	go func() {
		defer wg.Done()
		total, countErr = functionCollection.CountDocuments(ctx, bson.M{})
	}()

	// Fetch ALL paginated functions
	go func() {
		defer wg.Done()
		opts := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "created_at", Value: -1}})

		cursor, err := functionCollection.Find(ctx, bson.M{}, opts)
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

	// Apply privacy per function
	for i, fn := range allFunctions {
		allFunctions[i] = applyFunctionPrivacy(fn, userId)
	}

	// Cache processed data async
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

// -------------------- GET ONE FUNCTION (with privacy) --------------------
func GetOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid function ID"})
		return
	}

	cacheKey := fmt.Sprintf("function:%s", mongoId.Hex())
	var oneFunc models.Function

	// Try Redis
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			if jerr := json.Unmarshal([]byte(cached), &oneFunc); jerr == nil {
				fn := applyFunctionPrivacy(oneFunc, userId)
				c.JSON(200, gin.H{"msg": "Function from Redis✅", "function": fn, "source": "redis"})
				return
			}
		}
	}

	// Fetch without userId filter (ID se)
	if err := functionCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneFunc); err != nil {
		c.JSON(404, gin.H{"msg": "No function found❌"})
		return
	}

	// Cache raw event
	go func() {
		rctx := context.Background()
		if utils.RedisClient != nil {
			if b, err := json.Marshal(oneFunc); err == nil {
				_ = utils.RedisClient.Set(rctx, cacheKey, b, 60*time.Second).Err()
			}
		}
	}()

	fn := applyFunctionPrivacy(oneFunc, userId)

	c.JSON(200, gin.H{"msg": "Function from DB✅", "function": fn, "source": "db"})
}

// -------------------- EDIT FUNCTION (owner only) --------------------
func EditFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	funcId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid function ID"})
		return
	}

	// Pehle function exist check karo
	var existing models.Function
	if err := functionCollection.FindOne(ctx, bson.M{"_id": funcId}).Decode(&existing); err != nil {
		c.JSON(404, gin.H{"msg": "Function not found ❌"})
		return
	}

	// Ab ownership check
	if existing.UserId.Hex() != userId.Hex() {
		c.JSON(403, gin.H{
			"msg":    "Access denied 🚫",
			"reason": "Only the function owner can update this function.",
		})
		return
	}

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	// Step 1: Upload image (if given)
	imageUrl := existing.ImageUrl
	if file, _ := c.FormFile("file"); file != nil {
		if path, err := utils.FileUpload(c); err == nil {
			imageUrl = path
		}
	}

	// Step 2: Build update map dynamically
	updateFields := bson.M{}
	if funcName != "" {
		updateFields["funcname"] = funcName
	}
	if funcType != "" {
		updateFields["functype"] = funcType
	}
	if funcDesc != "" {
		updateFields["funcdes"] = funcDesc
	}
	if isPublic != "" {
		updateFields["ispublic"] = isPublic
	}
	if status != "" {
		updateFields["status"] = status
	}
	if location != "" {
		updateFields["location"] = location
	}
	updateFields["imageUrl"] = imageUrl
	updateFields["updatedAt"] = time.Now()

	// Step 3: Update in DB
	if _, err := functionCollection.UpdateByID(ctx, funcId, bson.M{"$set": updateFields}); err != nil {
		c.JSON(500, gin.H{"msg": "DB update failed", "error": err.Error()})
		return
	}

	// Step 4: Clear Redis cache async
	go func(uid, fid string) {
		rctx := context.Background()
		if utils.RedisClient != nil {
			_ = utils.RedisClient.Del(rctx, fmt.Sprintf("function:%s", fid)).Err()
			iter := utils.RedisClient.Scan(rctx, 0, fmt.Sprintf("function_list:%s:*", uid), 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}(userId.Hex(), funcId.Hex())

	c.JSON(200, gin.H{
		"msg":           "Function Updated Successfully ✅",
		"updatedFields": updateFields,
	})
}

// -------------------- DELETE ONE FUNCTION (owner only) --------------------
func DeleteOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid function ID"})
		return
	}

	// Pehle function exist check karo
	var fn models.Function
	if err := functionCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&fn); err != nil {
		c.JSON(404, gin.H{"msg": "Function not found ❌"})
		return
	}

	// Ownership check
	if fn.UserId.Hex() != userId.Hex() {
		c.JSON(403, gin.H{
			"msg":    "Access denied 🚫",
			"reason": "Only the function owner can delete this function.",
		})
		return
	}

	// Actual delete
	res, err := functionCollection.DeleteOne(ctx, bson.M{"_id": mongoId})
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error while deleting"})
		return
	}
	if res.DeletedCount == 0 {
		c.JSON(500, gin.H{"msg": "Delete failed unexpectedly"})
		return
	}

	// Clear Redis cache
	go func(fid string) {
		rctx := context.Background()
		if utils.RedisClient != nil {
			_ = utils.RedisClient.Del(rctx, fmt.Sprintf("function:%s", fid)).Err()
			pattern := "function_list:*"
			iter := utils.RedisClient.Scan(rctx, 0, pattern, 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}(mongoId.Hex())

	c.JSON(200, gin.H{"msg": "One Function Deleted✅"})
}

// -------------------- DELETE ALL FUNCTIONS (only own) --------------------
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
				"function:*",
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

// -------------------- JOIN FUNCTION --------------------
func JoinFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	funcID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Function ID"})
		return
	}

	funcColl := utils.MongoClient.Database("Event_Booking").Collection("functions")
	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	var function models.Function
	var existingReq models.JoinRequest

	// Step 1: Check if function exists
	if err := funcColl.FindOne(ctx, bson.M{"_id": funcID}).Decode(&function); err != nil {
		c.JSON(404, gin.H{"msg": "Function not found"})
		return
	}

	// Step 2: Check if user already requested or joined
	err = joinColl.FindOne(ctx, bson.M{
		"functionId":  funcID,
		"requesterId": userId,
	}).Decode(&existingReq)
	if err == nil {
		c.JSON(400, gin.H{"msg": "Already requested or joined"})
		return
	}

	// Step 3: Prevent self-join
	if function.UserId == userId {
		c.JSON(400, gin.H{"msg": "You cannot join your own function"})
		return
	}

	// Step 4: Public function — instant join
	if strings.ToLower(function.IsPublic) == "public" || strings.ToLower(function.IsPublic) == "true" {
		go func() {
			fmt.Printf("✅ User %s joined public function '%s'\n", userId.Hex(), function.FuncName)
		}()
		c.JSON(200, gin.H{
			"msg":           "Joined function successfully 🎉",
			"autoJoin":      true,
			"functionId":    function.ID,
			"functionName":  function.FuncName,
			"functionOwner": function.UserId,
		})
		return
	}

	// Step 5: Private function — create join request
	newReq := models.JoinRequest{
		ID:          primitive.NewObjectID(),
		FunctionID:  &funcID,
		RequesterID: userId,
		OwnerID:     function.UserId,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	go func() {
		rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := joinColl.InsertOne(rctx, newReq); err != nil {
			fmt.Println("❌ Join request insert failed:", err)
		} else {
			fmt.Printf("📨 Join request created for function '%s' by user %s\n", function.FuncName, userId.Hex())
		}
	}()

	c.JSON(200, gin.H{
		"msg":          "Join request sent successfully 📨",
		"functionName": function.FuncName,
		"requestId":    newReq.ID,
		"requestType":  "private function",
	})
}

// -------------------- APPROVE FUNCTION JOIN REQUEST --------------------
func ApproveFunctionJoinRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	ownerId := c.MustGet("userId").(primitive.ObjectID)
	reqID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request ID"})
		return
	}

	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	var req models.JoinRequest
	if err := joinColl.FindOne(ctx, bson.M{"_id": reqID, "ownerId": ownerId}).Decode(&req); err != nil {
		c.JSON(404, gin.H{"msg": "No such request or unauthorized"})
		return
	}

	done := make(chan bool, 1)
	go func() {
		rctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		_, err := joinColl.UpdateByID(rctx, reqID, bson.M{"$set": bson.M{
			"status":     "accepted",
			"updated_at": time.Now(),
		}})
		done <- err == nil
	}()

	if success := <-done; !success {
		c.JSON(500, gin.H{"msg": "Failed to approve join request"})
		return
	}

	c.JSON(200, gin.H{
		"msg":         "Join request approved ✅",
		"status":      "accepted",
		"functionId":  req.FunctionID,
		"requesterId": req.RequesterID,
	})
}

// -------------------- REJECT FUNCTION JOIN REQUEST --------------------
func RejectFunctionJoinRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ownerId := c.MustGet("userId").(primitive.ObjectID)
	reqID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request ID"})
		return
	}

	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	res, err := joinColl.UpdateOne(ctx,
		bson.M{"_id": reqID, "ownerId": ownerId},
		bson.M{"$set": bson.M{"status": "rejected", "updated_at": time.Now()}},
	)
	if err != nil || res.MatchedCount == 0 {
		c.JSON(404, gin.H{"msg": "No request found or unauthorized"})
		return
	}

	c.JSON(200, gin.H{"msg": "Join request rejected ❌"})
}

// -------------------- VIEW PENDING FUNCTION JOIN REQUESTS --------------------
func ViewPendingFunctionRequests(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	ownerId := c.MustGet("userId").(primitive.ObjectID)
	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	cursor, err := joinColl.Find(ctx, bson.M{
		"ownerId":    ownerId,
		"status":     "pending",
		"functionId": bson.M{"$exists": true},
	})
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error"})
		return
	}

	var requests []models.JoinRequest
	if err := cursor.All(ctx, &requests); err != nil {
		c.JSON(500, gin.H{"msg": "Failed to read join requests"})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "Pending function join requests ✅",
		"data": requests,
	})
}
