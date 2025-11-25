package ai

import (
	"time"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
)

func GenerateFunctionDesc(c *gin.Context) {
	var input AiFunctionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	function := models.Function{
		FuncName:  input.FuncName,
		FuncType:  input.FuncType,
		FuncDesc:  input.FuncDesc,
		Location:  input.Location,
		ImageUrl:  "no-image",
		IsPublic:  "public",
		Status:    "Upcoming",
		CreatedAt: time.Now(),
	}

	desc, err := utils.GenerateFunctionDescription(function)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to generate description"})
		return
	}

	c.JSON(200, gin.H{
		"funcname":    function.FuncName,
		"functype":    function.FuncType,
		"location":    function.Location,
		"description": desc,
	})
}
