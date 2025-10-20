package main

import (
	"github.com/AbdulRahman-04/07EvenetManagement-API/server/utils"
	"github.com/gin-gonic/gin"
)

func main(){

	// Mongodb & redis import 
	utils.DbConnect()
    utils.RedisConnect()

	router := gin.Default()

	router.GET("/", func (c*gin.Context)  {
		c.JSON(200, gin.H{
			"msg": "Hello world api",
		})
	})

	router.Run(":7575")
}