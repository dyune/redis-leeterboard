package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-leetcode-leaderboard/redis"
)

// About GET method
func About(context *gin.Context) {
	// Sends back JSON
	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"version":     "1.0.0",
			"application": "LeetCode leaderboard",
		})
}

// Register POST method
func Register(context *gin.Context) {
	var newUser redis.UserRequest

	// BindJSON binds the incoming request JSON to the defined type
	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(
			http.StatusBadRequest,
			err,
		)
	}

	user, err := redis.AddUser(newUser.Name, 0)
	if err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			err,
		)
	}

	// Sends back registered user
	context.IndentedJSON(
		http.StatusCreated,
		user,
	)
}
