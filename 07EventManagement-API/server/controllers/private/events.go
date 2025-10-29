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

var eventsCollection *mongo.Collection

func EventsCollect() {
	eventsCollection = utils.MongoClient.Database("Event_Booking").Collection("events")
}

func CreateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	// Read simple form fields
	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	// 1) Upload file first (if present) — do it synchronously to avoid races
	imageUrl := ""
	if file, _ := c.FormFile("file"); file != nil {
		if path, err := utils.FileUpload(c); err == nil {
			imageUrl = path
		} else {
			// don't fail the whole request if upload fails; keep empty image
			imageUrl = ""
		}
	}

	// 2) Convert attendance (simple, cheap operation)
	attendence := 0
	if eventAttendenceStr != "" {
		if v, err := strconv.Atoi(eventAttendenceStr); err == nil {
			attendence = v
		} else {
			c.JSON(400, gin.H{"msg": "Invalid attendance value", "error": err.Error()})
			return
		}
	}

	// 3) Build event struct and insert
	newEvent := models.Event{
		ID:               primitive.NewObjectID(),
		UserId:           userId,
		EventName:        eventName,
		EventtType:       eventType, // keep same field name as your model
		EventAttendence:  attendence,
		EventDescription: eventDes,
		IsPublic:         isPublic,
		Status:           status,
		Location:         location,
		ImageUrl:         imageUrl,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if _, err := eventsCollection.InsertOne(ctx, newEvent); err != nil {
		c.JSON(500, gin.H{"msg": "Failed to insert event into DB", "error": err.Error()})
		return
	}

	// 4) Async Redis cache invalidation (fire-and-forget)
	go func(uid string) {
		rctx := context.Background()
		if utils.RedisClient != nil {
			pattern := fmt.Sprintf("events:%s:*", uid)
			iter := utils.RedisClient.Scan(rctx, 0, pattern, 0).Iterator()
			for iter.Next(rctx) {
				_ = utils.RedisClient.Del(rctx, iter.Val()).Err()
			}
		}
	}(userId.Hex())

	c.JSON(200, gin.H{
		"msg":          "New Event Created ✨",
		"eventDetails": newEvent,
	})
}


// -------------------- GET ALL EVENTS --------------------
func GetAllEvents(c *gin.Context) {
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

	cacheKey := fmt.Sprintf("events:%s:%d:%d", userId.Hex(), page, limit)

	// Try from Redis cache first
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil && cached != "" {
			var cachedResponse struct {
				Msg     string         `json:"msg"`
				Events  []models.Event `json:"events"`
				Page    int            `json:"page"`
				Limit   int            `json:"limit"`
				Total   int64          `json:"total"`
				HasNext bool           `json:"hasNext"`
				HasPrev bool           `json:"hasPrev"`
				Source  string         `json:"source"`
			}
			if jsonErr := json.Unmarshal([]byte(cached), &cachedResponse); jsonErr == nil {
				cachedResponse.Source = "redis"
				c.JSON(200, cachedResponse)
				return
			}
		}
	}

	var (
		total     int64
		allEvents []models.Event
		countErr  error
		findErr   error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		total, countErr = eventsCollection.CountDocuments(ctx, bson.M{"userId": userId})
	}()

	go func() {
		defer wg.Done()
		opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "createdAt", Value: -1}})
		cursor, err := eventsCollection.Find(ctx, bson.M{"userId": userId}, opts)
		if err == nil {
			defer cursor.Close(ctx)
			findErr = cursor.All(ctx, &allEvents)
		} else {
			findErr = err
		}
	}()

	wg.Wait()
	if countErr != nil || findErr != nil {
		c.JSON(500, gin.H{"msg": "Failed to fetch events"})
		return
	}

	response := struct {
		Msg     string         `json:"msg"`
		Events  []models.Event `json:"events"`
		Page    int            `json:"page"`
		Limit   int            `json:"limit"`
		Total   int64          `json:"total"`
		HasNext bool           `json:"hasNext"`
		HasPrev bool           `json:"hasPrev"`
		Source  string         `json:"source"`
	}{
		Msg:     "All Events Are here✨",
		Events:  allEvents,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasNext: int64(skip+limit) < total,
		HasPrev: page > 1,
		Source:  "db",
	}

	// Cache response safely in background
	go func(resp any) {
		if utils.RedisClient != nil {
			cacheCtx := context.Background()
			dataBytes, _ := json.Marshal(resp)
			_ = utils.RedisClient.Set(cacheCtx, cacheKey, dataBytes, 60*time.Second).Err()
		}
	}(response)

	c.JSON(200, response)
}

// -------------------- GET ONE EVENT --------------------
func GetOneEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	cacheKey := fmt.Sprintf("event:%s:%s", userId.Hex(), mongoId.Hex())
	var oneEvent models.Event

	// Try Redis first
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(context.Background(), cacheKey).Result(); err == nil && cached != "" {
			if err := json.Unmarshal([]byte(cached), &oneEvent); err == nil {
				c.JSON(200, gin.H{"msg": "Event from Redis✅", "event": oneEvent, "source": "redis"})
				return
			}
		}
	}

	// Fallback to MongoDB
	if err := eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneEvent); err != nil {
		c.JSON(404, gin.H{"msg": "No event found❌"})
		return
	}

	// Cache the event
	go func(evt models.Event) {
		if utils.RedisClient != nil {
			cacheCtx := context.Background()
			dataBytes, _ := json.Marshal(evt)
			_ = utils.RedisClient.Set(cacheCtx, cacheKey, dataBytes, 60*time.Second).Err()
		}
	}(oneEvent)

	c.JSON(200, gin.H{"msg": "Event from DB✅", "event": oneEvent, "source": "db"})
}
func EditEventApi(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	// Fetch existing event
	var oldEvent models.Event
	if err := eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oldEvent); err != nil {
		c.JSON(404, gin.H{"msg": "No event found to update"})
		return
	}

	// Read form fields
	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	// 1) Upload new image (if provided). If upload fails, retain old image.
	imageUrl := oldEvent.ImageUrl
	if file, _ := c.FormFile("file"); file != nil {
		if path, err := utils.FileUpload(c); err == nil {
			imageUrl = path
		} else {
			// keep existing image if upload fails; do not abort update
			imageUrl = oldEvent.ImageUrl
		}
	}

	// 2) Parse attendance if provided
	var attendence int
	if eventAttendenceStr != "" {
		if v, err := strconv.Atoi(eventAttendenceStr); err == nil {
			attendence = v
		} else {
			c.JSON(400, gin.H{"msg": "Invalid attendance value", "error": err.Error()})
			return
		}
	}

	// 3) Build dynamic update map (only include provided fields)
	setMap := bson.M{}
	if eventName != "" {
		setMap["eventname"] = eventName
	}
	if eventType != "" {
		setMap["eventtype"] = eventType
	}
	if eventAttendenceStr != "" {
		setMap["attendence"] = attendence
	}
	if eventDes != "" {
		setMap["eventdesc"] = eventDes
	}
	if isPublic != "" {
		setMap["ispublic"] = isPublic
	}
	if status != "" {
		setMap["status"] = status
	}
	if location != "" {
		setMap["location"] = location
	}
	// imageUrl and updatedAt always set (imageUrl may be existing one)
	setMap["imageUrl"] = imageUrl
	setMap["updatedAt"] = time.Now()

	// 4) Execute update
	if _, err := eventsCollection.UpdateByID(ctx, mongoId, bson.M{"$set": setMap}); err != nil {
		c.JSON(500, gin.H{"msg": "Failed to update event", "error": err.Error()})
		return
	}

	// 5) Async: clear Redis cache keys related to this user/event
	go func(uid, fid string) {
		if utils.RedisClient != nil {
			cacheCtx := context.Background()
			_ = utils.RedisClient.Del(cacheCtx, fmt.Sprintf("event:%s:%s", uid, fid)).Err()

			pattern := fmt.Sprintf("events:%s:*", uid)
			iter := utils.RedisClient.Scan(cacheCtx, 0, pattern, 0).Iterator()
			for iter.Next(cacheCtx) {
				_ = utils.RedisClient.Del(cacheCtx, iter.Val()).Err()
			}
		}
	}(userId.Hex(), mongoId.Hex())

	c.JSON(200, gin.H{
		"msg":           "Event Updated Successfully ✅",
		"updatedFields": setMap,
	})
}


// -------------------- DELETE ONE EVENT --------------------
func DeleteOneEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	if _, err := eventsCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId}); err != nil {
		c.JSON(404, gin.H{"msg": "No Event Found or userId mismatch"})
		return
	}

	// Invalidate Redis cache
	go func() {
		if utils.RedisClient != nil {
			cacheCtx := context.Background()
			cacheKey := fmt.Sprintf("event:%s:%s", userId.Hex(), mongoId.Hex())
			_ = utils.RedisClient.Del(cacheCtx, cacheKey).Err()
		}
	}()

	c.JSON(200, gin.H{"msg": "One Event is deleted✅"})
}

// -------------------- DELETE ALL EVENTS --------------------
func DeleteAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	if _, err := eventsCollection.DeleteMany(ctx, bson.M{"userId": userId}); err != nil {
		c.JSON(500, gin.H{"msg": "DB error"})
		return
	}

	go func() {
		if utils.RedisClient != nil {
			cacheCtx := context.Background()
			pattern := fmt.Sprintf("event:%s:*", userId.Hex())
			iter := utils.RedisClient.Scan(cacheCtx, 0, pattern, 0).Iterator()
			for iter.Next(cacheCtx) {
				_ = utils.RedisClient.Del(cacheCtx, iter.Val()).Err()
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "All Events Deleted✅"})
}
