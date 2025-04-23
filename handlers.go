package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GET
func about(context *gin.Context) {
	// Sends back JSON
	context.IndentedJSON(
		http.StatusOK,
		gin.H{
			"version":     "1.0.0",
			"application": "LeetCode leaderboard",
		})
}

// POST
func register(context *gin.Context) {
	var newUser userRequest

	// BindJSON binds the incoming request JSON to the defined type
	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(
			http.StatusBadRequest,
			err,
		)
	}

	user, err := addUser(newUser.Name, 0)
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
