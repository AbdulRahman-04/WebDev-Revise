package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

	// pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	// ---------------- FETCH ALL EVENTS (NO FILTER) ----------------
	total, err := eventsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"msg": "Failed counting events"})
		return
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := eventsCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(500, gin.H{"msg": "Failed fetching events"})
		return
	}
	defer cursor.Close(ctx)

	var events []models.Event
	_ = cursor.All(ctx, &events)

	// ---------------- PRIVACY LOGIC ----------------
	for i, evt := range events {

		// normalize ispublic values ("public", "true" = treated as public)
		val := strings.ToLower(evt.IsPublic)
		isPublic := val == "public" || val == "true"
		isOwner := evt.UserId.Hex() == userId.Hex()

		// public → show full
		if isPublic {
			continue
		}

		// private but owner → show full
		if isOwner {
			continue
		}

		// private & NOT owner → show limited view
		events[i].Location = ""
		events[i].ImageUrl = ""
		events[i].EventAttendence = 0
		events[i].Status = "Private"
		events[i].EventDescription = "This is a private event. Request to join to see full details."
	}

	// ---------------- RESPONSE ----------------
	c.JSON(200, gin.H{
		"msg":     "Events fetched successfully ✨",
		"events":  events,
		"page":    page,
		"limit":   limit,
		"total":   total,
		"hasNext": int64(skip+limit) < total,
		"hasPrev": page > 1,
	})
}

// -------------------- GET ONE EVENT --------------------
func GetOneEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	mongoId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid event ID"})
		return
	}

	// --- Try Redis Cache First ---
	cacheKey := fmt.Sprintf("event:%s", mongoId.Hex())
	var event models.Event

	if utils.RedisClient != nil {
		if cached, _ := utils.RedisClient.Get(ctx, cacheKey).Result(); cached != "" {
			_ = json.Unmarshal([]byte(cached), &event)
			returnResponseWithPrivacy(c, event, userId, "redis")
			return
		}
	}

	// ---- Fetch Event without ownership restriction ----
	err = eventsCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&event)
	if err != nil {
		c.JSON(404, gin.H{"msg": "Event not found ❌"})
		return
	}

	// ---- Save in cache ---
	go func(evt models.Event) {
		if utils.RedisClient != nil {
			data, _ := json.Marshal(evt)
			utils.RedisClient.Set(context.Background(), cacheKey, data, 60*time.Second)
		}
	}(event)

	// ---- Final Output Based on Privacy ----
	returnResponseWithPrivacy(c, event, userId, "db")
}

func returnResponseWithPrivacy(c *gin.Context, evt models.Event, userId primitive.ObjectID, source string) {

	isPublic := strings.ToLower(evt.IsPublic) == "public" || strings.ToLower(evt.IsPublic) == "true"
	isOwner := evt.UserId.Hex() == userId.Hex()

	// If private and NOT owner -> mask
	if !isPublic && !isOwner {
		evt.Location = ""
		evt.ImageUrl = ""
		evt.Status = "Private"
		evt.EventAttendence = 0
		evt.EventDescription = "This is a private event. Request to join for full details."
	}

	c.JSON(200, gin.H{
		"msg":    "Event fetched successfully",
		"source": source,
		"event":  evt,
	})
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
		c.JSON(403, gin.H{
			"msg":    "Access denied 🚫",
			"reason": "You are not the owner of this event, so you cannot update it.",
		})
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
		c.JSON(400, gin.H{"msg": "Invalid event ID"})
		return
	}

	// TRY deleting event owned by this user
	result, err := eventsCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(500, gin.H{"msg": "Database delete error"})
		return
	}

	// ---- IMPORTANT FIX: Detect unauthorized delete attempt ----
	if result.DeletedCount == 0 {
		c.JSON(403, gin.H{
			"msg":    "Access denied 🚫",
			"reason": "Only the event owner can delete this event.",
		})
		return
	}

	// Clear Redis cache async
	go func() {
		if utils.RedisClient != nil {
			cacheCtx := context.Background()
			cacheKey := fmt.Sprintf("event:%s:%s", userId.Hex(), mongoId.Hex())
			_ = utils.RedisClient.Del(cacheCtx, cacheKey).Err()
		}
	}()

	c.JSON(200, gin.H{"msg": "Event deleted successfully ✅"})
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

// -------------------- JOIN EVENT (Optimized Concurrent Version) --------------------
func JoinEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	eventID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Event ID"})
		return
	}

	eventColl := utils.MongoClient.Database("Event_Booking").Collection("events")
	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	var (
		event   models.Event
		exists  models.JoinRequest
		foundCh = make(chan bool, 2)
	)

	// Concurrently check both event existence and previous request
	go func() {
		if err := eventColl.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event); err == nil {
			foundCh <- true
		} else {
			foundCh <- false
		}
	}()

	go func() {
		err := joinColl.FindOne(ctx, bson.M{"eventId": eventID, "requesterId": userId}).Decode(&exists)
		if err == nil {
			foundCh <- true
		} else {
			foundCh <- false
		}
	}()

	valid, already := <-foundCh, <-foundCh
	close(foundCh)

	if !valid {
		c.JSON(404, gin.H{"msg": "Event not found"})
		return
	}
	if already {
		c.JSON(400, gin.H{"msg": "Already requested or joined"})
		return
	}
	if event.UserId == userId {
		c.JSON(400, gin.H{"msg": "You cannot join your own event"})
		return
	}

	// ✅ Public event — instant join (no DB write)
	if event.IsPublic == "public" {
		go func() {
			fmt.Printf("User %s joined event '%s'\n", userId.Hex(), event.EventName)
		}()
		c.JSON(200, gin.H{
			"msg":       "You joined successfully 🎉",
			"autoJoin":  true,
			"eventName": event.EventName,
		})
		return
	}

	// 🔒 Private event — insert join request asynchronously
	newReq := models.JoinRequest{
		ID:          primitive.NewObjectID(),
		EventID:     eventID,
		RequesterID: userId,
		OwnerID:     event.UserId,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	go func() {
		rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := joinColl.InsertOne(rctx, newReq); err != nil {
			fmt.Println("❌ Join request insert failed:", err)
		}
	}()

	c.JSON(200, gin.H{
		"msg":        "Join request sent successfully 📨",
		"eventName":  event.EventName,
		"requestId":  newReq.ID,
		"requestFor": "private event",
	})
}

// -------------------- APPROVE JOIN REQUEST (Optimized) --------------------
func ApproveJoinRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	ownerId := c.MustGet("userId").(primitive.ObjectID)
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request ID"})
		return
	}

	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	var req models.JoinRequest
	if err := joinColl.FindOne(ctx, bson.M{"_id": requestID, "ownerId": ownerId}).Decode(&req); err != nil {
		c.JSON(404, gin.H{"msg": "No such request or unauthorized"})
		return
	}

	// Concurrent update
	done := make(chan bool, 1)
	go func() {
		rctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		_, err := joinColl.UpdateByID(rctx, requestID, bson.M{"$set": bson.M{
			"status":     "accepted",
			"updated_at": time.Now(),
		}})
		if err != nil {
			fmt.Println("❌ Update error:", err)
			done <- false
		} else {
			done <- true
		}
	}()

	if success := <-done; !success {
		c.JSON(500, gin.H{"msg": "Failed to approve request"})
		return
	}

	c.JSON(200, gin.H{
		"msg":     "Join request approved ✅",
		"status":  "accepted",
		"eventId": req.EventID,
		"userId":  req.RequesterID,
	})
}

func RejectJoinRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ownerId := c.MustGet("userId").(primitive.ObjectID)
	requestID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid Request ID"})
		return
	}

	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")
	update := bson.M{"$set": bson.M{"status": "rejected", "updated_at": time.Now()}}
	res, err := joinColl.UpdateOne(ctx, bson.M{"_id": requestID, "ownerId": ownerId}, update)
	if err != nil || res.MatchedCount == 0 {
		c.JSON(404, gin.H{"msg": "No such request or unauthorized"})
		return
	}

	c.JSON(200, gin.H{"msg": "Join request rejected ❌"})
}

func ViewPendingRequests(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ownerId := c.MustGet("userId").(primitive.ObjectID)
	joinColl := utils.MongoClient.Database("Event_Booking").Collection("join_requests")

	cursor, err := joinColl.Find(ctx, bson.M{"ownerId": ownerId, "status": "pending"})
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB error"})
		return
	}
	var requests []models.JoinRequest
	_ = cursor.All(ctx, &requests)

	c.JSON(200, gin.H{"msg": "All pending join requests", "data": requests})
}
