package api

import (
	"github.com/gin-gonic/gin"
	"redis-leetcode-leaderboard/service"
)

func SetupRoutes(router *gin.RouterGroup) {

	// GET: Information on application
	router.GET("/ping", service.About)

	// POST: Register a user
	router.POST("/users", service.Register)

	// GET: Get information on a user
	router.GET("/users/:id", service.GetUserData)

	// POST: Update a user's score
	router.POST("/users/:id/score", service.UpdateUserPoints)

	// DELETE: Delete a user from the leaderboard
	router.DELETE("/users/:id", service.DeleteUser)

	// GET: Get the whole leaderboard
	router.GET("/leaderboard", service.GetLeaderboard)

}
