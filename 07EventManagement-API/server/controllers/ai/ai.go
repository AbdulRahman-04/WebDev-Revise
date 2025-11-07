package ai

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
)

// ✅ Generate Event Description (AI + all model data)
func GenerateDescription(c *gin.Context) {
	var event models.Event

	// Bind JSON input
	if err := c.BindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call AI util function (uses all event model data)
	desc, err := utils.GenerateEventDescription(event)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to generate event description"})
		return
	}

	// Send response
	c.JSON(200, gin.H{
		"eventname":   event.EventName,
		"eventtype":   event.EventtType,
		"attendence":  event.EventAttendence,
		"location":    event.Location,
		"status":      event.Status,
		"ispublic":    event.IsPublic,
		"description": desc,
	})
}

// ✅ Generate Function Description (AI + model data)
func GenerateFunctionDesc(c *gin.Context) {
	var function models.Function

	// Bind JSON input
	if err := c.BindJSON(&function); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call AI util function (uses all function model data)
	desc, err := utils.GenerateFunctionDescription(function)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to generate function description"})
		return
	}

	// Send response
	c.JSON(200, gin.H{
		"funcname":    function.FuncName,
		"functype":    function.FuncType,
		"location":    function.Location,
		"status":      function.Status,
		"ispublic":    function.IsPublic,
		"description": desc,
	})
}
