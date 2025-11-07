package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RecommendAI(c *gin.Context) {
	var events []models.Event
	var functions []models.Function

	eventColl := utils.MongoClient.Database("Event_Booking").Collection("events")
	funcColl := utils.MongoClient.Database("Event_Booking").Collection("functions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventCursor, err := eventColl.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch events"})
		return
	}
	_ = eventCursor.All(ctx, &events)

	funcCursor, err := funcColl.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch functions"})
		return
	}
	_ = funcCursor.All(ctx, &functions)

	if len(events) == 0 && len(functions) == 0 {
		c.JSON(200, gin.H{"message": "No events or functions found in DB"})
		return
	}

	// Prepare data text for AI
	var eventList, funcList string
	for _, e := range events {
		eventList += fmt.Sprintf("- %s (%s) at %s\n", e.EventName, e.EventtType, e.Location)
	}
	for _, f := range functions {
		funcList += fmt.Sprintf("- %s (%s) at %s\n", f.FuncName, f.FuncType, f.Location)
	}

	// 🧠 Prompt
	prompt := fmt.Sprintf(`
You are an intelligent event recommendation system.

Here are all events and functions in the database:
EVENTS:
%s
FUNCTIONS:
%s

Recommend 5 events and 5 functions (only from above list).
Respond in pure JSON (no markdown, no code):
{
  "recommendedEvents": ["Event1", "Event2", "Event3", "Event4", "Event5"],
  "recommendedFunctions": ["Func1", "Func2", "Func3", "Func4", "Func5"]
}
`, eventList, funcList)

	result, err := utils.GenerateAIResponse(prompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to generate recommendations"})
		return
	}

	// Extract names
	var recEvents []string
	var recFuncs []string

	if evs, ok := result["recommendedEvents"].([]interface{}); ok {
		for _, e := range evs {
			if str, ok := e.(string); ok && str != "" {
				recEvents = append(recEvents, strings.TrimSpace(str))
			}
		}
	}
	if funcs, ok := result["recommendedFunctions"].([]interface{}); ok {
		for _, f := range funcs {
			if str, ok := f.(string); ok && str != "" {
				recFuncs = append(recFuncs, strings.TrimSpace(str))
			}
		}
	}

	// ✅ Use maps to avoid duplicates
	uniqueEventIDs := make(map[string]bool)
	uniqueFuncIDs := make(map[string]bool)

	var finalEvents []models.Event
	for _, name := range recEvents {
		clean := strings.Split(name, "(")[0]
		clean = strings.Split(clean, "at")[0]
		clean = strings.TrimSpace(clean)

		var evt models.Event
		err := eventColl.FindOne(ctx, bson.M{
			"eventname": bson.M{"$regex": clean, "$options": "i"},
		}).Decode(&evt)

		if err == nil && evt.ID != primitive.NilObjectID {
			if !uniqueEventIDs[evt.ID.Hex()] {
				finalEvents = append(finalEvents, evt)
				uniqueEventIDs[evt.ID.Hex()] = true
			}
		}
	}

	var finalFuncs []models.Function
	for _, name := range recFuncs {
		clean := strings.Split(name, "(")[0]
		clean = strings.Split(clean, "at")[0]
		clean = strings.TrimSpace(clean)

		var fn models.Function
		err := funcColl.FindOne(ctx, bson.M{
			"funcname": bson.M{"$regex": clean, "$options": "i"},
		}).Decode(&fn)

		if err == nil && fn.ID != primitive.NilObjectID {
			if !uniqueFuncIDs[fn.ID.Hex()] {
				finalFuncs = append(finalFuncs, fn)
				uniqueFuncIDs[fn.ID.Hex()] = true
			}
		}
	}

	// ✅ Return clean response
	c.JSON(200, gin.H{
		"type": "recommendations",
		"data": gin.H{
			"recommendedEvents":    finalEvents,
			"recommendedFunctions": finalFuncs,
		},
	})
}
