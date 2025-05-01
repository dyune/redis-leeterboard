package api

import (
	"github.com/gin-gonic/gin"
	"redis-leetcode-leaderboard/service"
)

func SetupRoutes(router *gin.RouterGroup) {

	// GET: Information on application
	router.GET("/", service.About)

	// POST: Register a user
	router.POST("/users", service.Register)

	// GET: Get information on a user
	router.GET("/users/:id", service.GetUserData)

}
