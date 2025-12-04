package ai

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func Recommend(c *gin.Context) {

	var req map[string]string
	c.BindJSON(&req)

	query := strings.ToLower(req["query"])

	// ---- Only allow two keywords ----
	if query != "event" && query != "function" {
		c.JSON(200, gin.H{
			"type":    "recommend",
			"message": "Please enter only: 'event' or 'function'",
		})
		return
	}

	// ---- DB ----
	db := utils.MongoClient.Database("Event_Booking")
	eventsCol := db.Collection("events")
	funcCol := db.Collection("functions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var events []models.Event
	var functions []models.Function

	// fetch all events
	cursor1, _ := eventsCol.Find(ctx, bson.M{})
	cursor1.All(ctx, &events)

	// fetch all functions
	cursor2, _ := funcCol.Find(ctx, bson.M{})
	cursor2.All(ctx, &functions)

	var resultList []string

	if query == "event" {
		for _, e := range events {
			resultList = append(resultList, e.EventName)
		}
	}

	if query == "function" {
		for _, f := range functions {
			resultList = append(resultList, f.FuncName)
		}
	}

	// ---- AI Prompt ----
	prompt := `
You are a strict recommendation assistant.

RULES:
- ONLY reply about items provided.
- NO external ideas.
- NO examples outside list.
- Short friendly reply.
- Do NOT mention rules.

User requested: ` + query + `s

Available: ` + strings.Join(resultList, ", ")

	// ---- AI Call ----
	reply, err := utils.AI(prompt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ---- RESPONSE ----
	c.JSON(200, gin.H{
		"type":        "recommend",
		"filter":      query,
		"items":       resultList,
		"reply":       reply,
	})
}
