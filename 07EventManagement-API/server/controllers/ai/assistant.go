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

// ==========================================================
// 🧠 FINAL OPTION-A: SMART CHAT + TOP 3 RECOMMENDATIONS
// ==========================================================
func Assistant(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	db := utils.MongoClient.Database("Event_Booking")

	eventColl := db.Collection("events")
	funcColl := db.Collection("functions")
	joinColl := db.Collection("join_requests")

	// ----------------------------
	// 1. Read Prompt
	// ----------------------------
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	query := strings.ToLower(req.Prompt)

	// Intent detection
	onlyFunctions := strings.Contains(query, "function") ||
		strings.Contains(query, "shaadi") ||
		strings.Contains(query, "valima") ||
		strings.Contains(query, "aqeeqa") ||
		strings.Contains(query, "birthday")

	onlyEvents := strings.Contains(query, "event") ||
		strings.Contains(query, "party") ||
		strings.Contains(query, "meetup") ||
		strings.Contains(query, "conference")

	// ----------------------------
	// 2. Fetch Events
	// ----------------------------
	var events []models.Event
	if !onlyFunctions {
		cur, _ := eventColl.Find(ctx, bson.M{})
		_ = cur.All(ctx, &events)
	}

	// ----------------------------
	// 3. Fetch Functions
	// ----------------------------
	var functions []models.Function
	if !onlyEvents {
		cur, _ := funcColl.Find(ctx, bson.M{})
		_ = cur.All(ctx, &functions)
	}

	// ----------------------------
	// 4. User Joined History
	// ----------------------------
	var joined []models.JoinRequest
	cur, _ := joinColl.Find(ctx, bson.M{"requesterId": userId, "status": "accepted"})
	_ = cur.All(ctx, &joined)

	var joinedEvents []string
	var joinedFuncs []string

	for _, j := range joined {
		if !j.EventID.IsZero() {
			var e models.Event
			if err := eventColl.FindOne(ctx, bson.M{"_id": j.EventID}).Decode(&e); err == nil {
				joinedEvents = append(joinedEvents, e.EventName)
			}
		}
		if j.FunctionID != nil {
			var f models.Function
			if err := funcColl.FindOne(ctx, bson.M{"_id": *j.FunctionID}).Decode(&f); err == nil {
				joinedFuncs = append(joinedFuncs, f.FuncName)
			}
		}
	}

	// ----------------------------
	// 5. Prepare JSON lists (clean)
	// ----------------------------
	eventList := "["
	for _, e := range events {
		eventList += fmt.Sprintf(`{"name":"%s","type":"%s","location":"%s","visibility":"%s","status":"%s"},`,
			e.EventName, e.EventtType, e.Location, e.IsPublic, e.Status)
	}
	eventList += "]"

	funcList := "["
	for _, f := range functions {
		funcList += fmt.Sprintf(`{"name":"%s","type":"%s","location":"%s","visibility":"%s","status":"%s"},`,
			f.FuncName, f.FuncType, f.Location, f.IsPublic, f.Status)
	}
	funcList += "]"

	// ----------------------------
	// 6. FINAL PROMPT (Smart)
	// ----------------------------
	aiPrompt := fmt.Sprintf(`
You are a professional event/function assistant AI.

USER QUERY:
"%s"

USER HISTORY:
Joined Events: %v
Joined Functions: %v

AVAILABLE EVENTS (JSON):
%s

AVAILABLE FUNCTIONS (JSON):
%s

IMPORTANT RULES:

1. If query is about FUNCTIONS → recommend ONLY functions.
2. If query is about EVENTS → recommend ONLY events.
3. If query is general → recommend BOTH.

🔥 MOST IMPORTANT:
- ALWAYS return ALL relevant functions/events from the database.
- DO NOT return just one best match.
- If 2 functions exist in the DB and both are relevant → return BOTH.
- Output max 10 items. Never just 1.

OUTPUT MUST be pure JSON ONLY:

{
  "answer": "...",
  "recommendedEvents": [{"name": "...", "reason": "...", "score": 0-100}],
  "recommendedFunctions": [{"name": "...", "reason": "...", "score": 0-100}]
}
`, req.Prompt, joinedEvents, joinedFuncs, eventList, funcList)


	// ----------------------------
	// 7. AI Call
	// ----------------------------
	aiResp, err := utils.GenerateAIResponse(aiPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed"})
		return
	}

	// ----------------------------
	// 8. Final Response
	// ----------------------------
	c.JSON(200, gin.H{
		"type":     "ai-assistant",
		"response": aiResp,
	})
}
