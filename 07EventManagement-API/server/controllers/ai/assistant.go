package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 🚀 Industry-grade personalized AI recommendation engine
func RecommendAI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	db := utils.MongoClient.Database("Event_Booking")

	eventColl := db.Collection("events")
	funcColl := db.Collection("functions")
	joinColl := db.Collection("join_requests")

	var (
		allEvents     []models.Event
		allFunctions  []models.Function
		userJoinReqs  []models.JoinRequest
		wg            sync.WaitGroup
		eventErr      error
		funcErr       error
		joinErr       error
	)

	// 🧵 Concurrent fetching for performance
	wg.Add(3)
	go func() {
		defer wg.Done()
		// fetch both public and private for better recommendation context
		cursor, err := eventColl.Find(ctx, bson.M{
			"$or": []bson.M{
				{"ispublic": "public"},
				{"ispublic": "private"},
			},
		})
		if err != nil {
			eventErr = err
			return
		}
		defer cursor.Close(ctx)
		_ = cursor.All(ctx, &allEvents)
	}()
	go func() {
		defer wg.Done()
		cursor, err := funcColl.Find(ctx, bson.M{
			"$or": []bson.M{
				{"ispublic": "public"},
				{"ispublic": "private"},
			},
		})
		if err != nil {
			funcErr = err
			return
		}
		defer cursor.Close(ctx)
		_ = cursor.All(ctx, &allFunctions)
	}()
	go func() {
		defer wg.Done()
		cursor, err := joinColl.Find(ctx, bson.M{"requesterId": userId, "status": "accepted"})
		if err != nil {
			joinErr = err
			return
		}
		defer cursor.Close(ctx)
		_ = cursor.All(ctx, &userJoinReqs)
	}()
	wg.Wait()

	if eventErr != nil || funcErr != nil || joinErr != nil {
		c.JSON(500, gin.H{"error": "DB fetch failed"})
		return
	}

	// ⚙️ Extract joined event and function names for personalization
	var joinedEventNames, joinedFuncNames []string
	for _, j := range userJoinReqs {
		if !j.EventID.IsZero() {
			var evt models.Event
			if err := eventColl.FindOne(ctx, bson.M{"_id": j.EventID}).Decode(&evt); err == nil {
				joinedEventNames = append(joinedEventNames, evt.EventName)
			}
		}
		if j.FunctionID != nil {
			var fn models.Function
			if err := funcColl.FindOne(ctx, bson.M{"_id": *j.FunctionID}).Decode(&fn); err == nil {
				joinedFuncNames = append(joinedFuncNames, fn.FuncName)
			}
		}
	}

	// 🧩 Create user profile context for AI
	userProfile := fmt.Sprintf(`
	User ID: %s
	Joined Events: %v
	Joined Functions: %v
	`, userId.Hex(), joinedEventNames, joinedFuncNames)

	// 🧾 Prepare dataset for AI
	var eventList, funcList string
	for _, e := range allEvents {
		eventList += fmt.Sprintf("- %s (%s) at %s [%s]\n", e.EventName, e.EventtType, e.Location, e.IsPublic)
	}
	for _, f := range allFunctions {
		funcList += fmt.Sprintf("- %s (%s) at %s [%s]\n", f.FuncName, f.FuncType, f.Location, f.IsPublic)
	}

	// 🧠 Smart prompt to AI
	prompt := fmt.Sprintf(`
	You are an advanced event recommendation AI.

	User Profile:
	%s

	Available Events:
	%s

	Available Functions:
	%s

	Recommend the top 5 events and 5 functions this user would most likely attend.
	Return pure JSON only, no markdown:
	{
	"recommendedEvents": [{"name": "...", "reason": "..."}],
	"recommendedFunctions": [{"name": "...", "reason": "..."}]
	}
	`, userProfile, eventList, funcList)

	// 🔮 AI response
	result, err := utils.GenerateAIResponse(prompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to generate recommendations"})
		return
	}

	// 🧾 Extract AI data safely
	var recEventsRaw, recFuncsRaw []map[string]interface{}
	if evs, ok := result["recommendedEvents"].([]interface{}); ok {
		for _, e := range evs {
			if m, ok := e.(map[string]interface{}); ok {
				recEventsRaw = append(recEventsRaw, m)
			}
		}
	}
	if funcs, ok := result["recommendedFunctions"].([]interface{}); ok {
		for _, f := range funcs {
			if m, ok := f.(map[string]interface{}); ok {
				recFuncsRaw = append(recFuncsRaw, m)
			}
		}
	}

	// 🧩 Build final detailed recommendations
	var finalEvents []gin.H
	for _, item := range recEventsRaw {
		name := strings.TrimSpace(fmt.Sprintf("%v", item["name"]))
		reason := strings.TrimSpace(fmt.Sprintf("%v", item["reason"]))
		var evt models.Event
		if err := eventColl.FindOne(ctx, bson.M{"eventname": bson.M{"$regex": name, "$options": "i"}}).Decode(&evt); err == nil {
			finalEvents = append(finalEvents, gin.H{
				"eventname": evt.EventName,
				"location":  evt.Location,
				"status":    evt.Status,
				"ispublic":  evt.IsPublic,
				"image":     evt.ImageUrl,
				"reason":    reason,
			})
		}
	}

	var finalFuncs []gin.H
	for _, item := range recFuncsRaw {
		name := strings.TrimSpace(fmt.Sprintf("%v", item["name"]))
		reason := strings.TrimSpace(fmt.Sprintf("%v", item["reason"]))
		var fn models.Function
		if err := funcColl.FindOne(ctx, bson.M{"funcname": bson.M{"$regex": name, "$options": "i"}}).Decode(&fn); err == nil {
			finalFuncs = append(finalFuncs, gin.H{
				"funcname": fn.FuncName,
				"location": fn.Location,
				"status":   fn.Status,
				"ispublic": fn.IsPublic,
				"image":    fn.ImageUrl,
				"reason":   reason,
			})
		}
	}

	// ✅ Final clean JSON response
	c.JSON(200, gin.H{
		"type": "personalized_recommendations",
		"user": userId.Hex(),
		"data": gin.H{
			"recommendedEvents":    finalEvents,
			"recommendedFunctions": finalFuncs,
		},
	})
}
