package ai

import (
	"time"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
)

func GenerateDescription(c *gin.Context) {
	var input AiEventInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{
		EventName:        input.EventName,
		EventtType:       input.EventType,
		EventAttendence:  input.Attendence,
		EventDescription: input.EventDesc,
		Location:         input.Location,
		ImageUrl:         "no-image",
		IsPublic:         "public",
		Status:           "Upcoming",
		CreatedAt:        time.Now(),
	}

	desc, err := utils.GenerateEventDescription(event)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to generate description"})
		return
	}

	c.JSON(200, gin.H{
		"eventname":   event.EventName,
		"eventtype":   event.EventtType,
		"attendence":  event.EventAttendence,
		"location":    event.Location,
		"description": desc,
	})
}
