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

// -------------------- CREATE EVENT --------------------
func CreateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	var (
		imageUrl    string
		uploadErr   error
		insertErr   error
		attendence  int
		convertErr  error
	)

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
		attendence, convertErr = strconv.Atoi(eventAttendenceStr)
	}()

	wg.Wait()
	if convertErr != nil {
		c.JSON(400, gin.H{"msg": "Conversion error"})
		return
	}

	newEvent := models.Event{
		ID:               primitive.NewObjectID(),
		UserId:           userId,
		EventName:        eventName,
		EventtType:       eventType,
		EventAttendence:  attendence,
		EventDescription: eventDes,
		IsPublic:         isPublic,
		Status:           status,
		Location:         location,
		ImageUrl:         imageUrl,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, insertErr = eventsCollection.InsertOne(ctx, newEvent)
	if insertErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	c.JSON(200, gin.H{"msg": "New Event Created✨", "eventDetails": newEvent})
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
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
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
		total      int64
		allEvents  []models.Event
		countErr   error
		findErr    error
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
		c.JSON(400, gin.H{"msg": "DB error"})
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

	go func() {
		if utils.RedisClient != nil {
			cacheResp := response
			cacheResp.Source = ""
			dataBytes, _ := json.Marshal(cacheResp)
			_ = utils.RedisClient.Set(ctx, cacheKey, dataBytes, 60*time.Second).Err()
		}
	}()

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
	var redisErr, dbErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if utils.RedisClient != nil {
			if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
				redisErr = json.Unmarshal([]byte(cached), &oneEvent)
			}
		}
	}()

	go func() {
		defer wg.Done()
		dbErr = eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneEvent)
	}()

	wg.Wait()

	if redisErr == nil {
		c.JSON(200, gin.H{"msg": "Event from Redis✅", "event": oneEvent, "source": "redis"})
		return
	}
	if dbErr != nil {
		c.JSON(404, gin.H{"msg": "No event found❌"})
		return
	}

	go func() {
		if utils.RedisClient != nil {
			dataBytes, _ := json.Marshal(oneEvent)
			_ = utils.RedisClient.Set(ctx, cacheKey, dataBytes, 60*time.Second).Err()
		}
	}()

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

	var oldEvent models.Event
	if err := eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oldEvent); err != nil {
		c.JSON(400, gin.H{"msg": "No event found to update"})
		return
	}

	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")

	var (
		imageUrl   string
		uploadErr  error
		attendence int
		convertErr error
		updateErr  error
	)

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
		attendence, convertErr = strconv.Atoi(eventAttendenceStr)
	}()

	wg.Wait()
	if convertErr != nil {
		c.JSON(400, gin.H{"msg": "Conversion error"})
		return
	}

	update := bson.M{"$set": bson.M{
		"eventname":  eventName,
		"eventtype":  eventType,
		"attendence": attendence,
		"eventdesc":  eventDes,
		"ispublic":   isPublic,
		"status":     status,
		"location":   location,
		"imageUrl":   imageUrl,
		"updated_at": time.Now(),
	}}

	_, updateErr = eventsCollection.UpdateByID(ctx, mongoId, update)
	if updateErr != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Event Updated Successfully!✅", "UpdatedEvent": oldEvent})
}
func DeleteOneEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	_, err = eventsCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "No Event Found or userId mismatch"})
		return
	}

	go func() {
		cacheKey := fmt.Sprintf("event:%s:%s", userId.Hex(), mongoId.Hex())
		if utils.RedisClient != nil {
			_ = utils.RedisClient.Del(ctx, cacheKey).Err()
		}
	}()

	c.JSON(200, gin.H{"msg": "One Event is deleted✅"})
}

func DeleteAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	_, err := eventsCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "DB error"})
		return
	}

	go func() {
		if utils.RedisClient != nil {
			pattern := fmt.Sprintf("event:%s:*", userId.Hex())
			iter := utils.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
			for iter.Next(ctx) {
				_ = utils.RedisClient.Del(ctx, iter.Val()).Err()
			}
		}
	}()

	c.JSON(200, gin.H{"msg": "All Events Deleted✅"})
}

