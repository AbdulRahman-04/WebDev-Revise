package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateDescription(c *gin.Context) {

	var req struct {
		ID string `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	db := utils.MongoClient.Database("Event_Booking")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ================================
	// 🔍 1) TRY FIND EVENT
	// ================================
	var event models.Event
	err = db.Collection("events").FindOne(ctx, bson.M{"_id": objID}).Decode(&event)

	if err == nil {
		// Create AI prompt using event fields
		prompt := fmt.Sprintf(`
Write a clean, short (max 90 words), friendly event description.

Event Name: %s
Type: %s
Expected Attendance: %d people
Visibility: %s
Status: %s
Location: %s

Rules:
- No emojis ❌
- No hashtags ❌
- No bullet points ❌
- Should sound natural and helpful.
`, event.EventName, event.EventtType, event.EventAttendence, event.IsPublic, event.Status, event.Location)

		response, aiErr := utils.AI(prompt)
		if aiErr != nil {
			c.JSON(500, gin.H{"error": aiErr.Error()})
			return
		}

		c.JSON(200, gin.H{
			"type":        "event",
			"id":          req.ID,
			"title":       event.EventName,
			"description": response,
		})
		return
	}

	// ================================
	// 🔍 2) TRY FIND FUNCTION
	// ================================
	var fn models.Function
	err = db.Collection("functions").FindOne(ctx, bson.M{"_id": objID}).Decode(&fn)

	if err == nil {
		prompt := fmt.Sprintf(`
Write a warm, elegant and short function description (max 80 words).

Function Name: %s
Type: %s
Visibility: %s
Status: %s
Location: %s

Tone: human, emotional, inviting. No emojis, no bullet points.
`, fn.FuncName, fn.FuncType, fn.IsPublic, fn.Status, fn.Location)

		response, aiErr := utils.AI(prompt)
		if aiErr != nil {
			c.JSON(500, gin.H{"error": aiErr.Error()})
			return
		}

		c.JSON(200, gin.H{
			"type":        "function",
			"id":          req.ID,
			"title":       fn.FuncName,
			"description": response,
		})
		return
	}

	// ❌ NO MATCH
	c.JSON(404, gin.H{"error": "No matching event or function found"})
}
