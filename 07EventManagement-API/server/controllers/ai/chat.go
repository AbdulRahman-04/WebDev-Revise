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

func Chat(c *gin.Context) {

	var req map[string]string
	c.BindJSON(&req)

	// ------------------ USER MESSAGE CLEAN ------------------
	rawMsg := strings.TrimSpace(req["message"])
	safeMsg := strings.ReplaceAll(rawMsg, `"`, `'`)
	lower := strings.ToLower(safeMsg)

	// ------------------ DB SETUP ------------------
	db := utils.MongoClient.Database("Event_Booking")
	eventsCol := db.Collection("events")
	funcCol := db.Collection("functions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ------------------ FETCH ALL FOR AI CONTEXT ------------------
	var events []models.Event
	cursor1, _ := eventsCol.Find(ctx, bson.M{})
	_ = cursor1.All(ctx, &events)

	var functions []models.Function
	cursor2, _ := funcCol.Find(ctx, bson.M{})
	_ = cursor2.All(ctx, &functions)

	// Convert names list
	eventNames := []string{}
	for _, e := range events {
		eventNames = append(eventNames, e.EventName)
	}

	funcNames := []string{}
	for _, f := range functions {
		funcNames = append(funcNames, f.FuncName)
	}

	// ------------------ INTENT FILTER (just for output JSON) ------------------
	isEvent := strings.Contains(lower, "event")
	isFunction := strings.Contains(lower, "function") ||
		strings.Contains(lower, "shaadi") ||
		strings.Contains(lower, "birthday") ||
		strings.Contains(lower, "valima") ||
		strings.Contains(lower, "mehendi") ||
		strings.Contains(lower, "aqeeqa")

	sendEvents := eventNames
	sendFuncs := funcNames

	if isEvent && !isFunction {
		sendFuncs = []string{}
	} else if isFunction && !isEvent {
		sendEvents = []string{}
	}

	// ------------------ AI PROMPT ------------------
	prompt := `
You are a friendly assistant inside an Event & Function booking app.

Talk only about:
✔️ Events
✔️ Functions
✔️ Joining / hosting inside the app

Never go off topic.

App Data:
Events: ` + strings.Join(eventNames, ", ") + `
Functions: ` + strings.Join(funcNames, ", ") + `

User: "` + safeMsg + `"

Reply short (2-3 sentences), friendly, helpful.
`

	// ------------------ AI CALL ------------------
	reply, err := utils.AI(prompt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ------------------ SEND RESPONSE ------------------
	c.JSON(200, gin.H{
		"type":      "chat",
		"reply":     reply,
		"events":    sendEvents,
		"functions": sendFuncs,
	})
}
