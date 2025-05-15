package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-leetcode-leaderboard/redis"
	"strconv"
	"time"
)

// About GET method
func About(context *gin.Context) {

	start := time.Now()

	response, err := redis.Rdb.Ping(context).Result()

	duration := time.Since(start)

	if err != nil {
		context.IndentedJSON(
			http.StatusServiceUnavailable,
			gin.H{
				"status":        "unhealthy",
				"error":         err,
				"response":      response,
				"redis_latency": fmt.Sprintf("%dms", duration.Milliseconds()),
			},
		)
	} else {
		context.IndentedJSON(
			http.StatusOK,
			gin.H{
				"status":        "healthy",
				"response":      response,
				"redis_latency": fmt.Sprintf("%dms", duration.Milliseconds()),
				"version":       "1.0.0",
				"application":   "LeetCode leaderboard",
			})
	}
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
		return
	}

	user, err := redis.AddUser(newUser.Name, 0.00)
	if err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	// Sends back registered user
	context.IndentedJSON(
		http.StatusCreated,
		user,
	)
}

// GetUserData GET method
func GetUserData(context *gin.Context) {

	id := context.Param("id")

	data, err := redis.GetUserData(id)
	if err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			err,
		)
	} else {
		// Send back user information
		context.IndentedJSON(
			http.StatusOK,
			data,
		)
	}
}

// UpdateUserPoints POST method
func UpdateUserPoints(context *gin.Context) {

	id := context.Param("id")
	var points redis.PointUpdateRequest

	err := context.BindJSON(&points)
	if err != nil {
		context.IndentedJSON(
			http.StatusBadRequest,
			err,
		)
	}

	if score, err := redis.IncreaseUserScore(id, points.Points); err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			err,
		)
	} else {
		context.IndentedJSON(
			http.StatusOK,
			gin.H{
				"newScore": score,
			},
		)
	}

}

// GetLeaderboard GET method
func GetLeaderboard(context *gin.Context) {

	leaderboard, err := redis.GetLeaderboard()

	if err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
	} else {
		context.IndentedJSON(
			http.StatusOK,
			leaderboard,
		)
	}
}

// DeleteUser DELETE method
func DeleteUser(context *gin.Context) {

	id := context.Param("id")

	err := redis.DeleteUser(id)

	if err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			err,
		)
	} else {
		context.IndentedJSON(
			http.StatusOK,
			gin.H{
				"Deleted-user-ID": id,
			},
		)
	}
}

// GetSubset GET method
func GetSubset(context *gin.Context) {
	offset, err := strconv.Atoi(context.DefaultQuery("offset", "1"))
	if err != nil {
		context.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": "`offset` must be a valid integer"})
		return
	}
	limit, err := strconv.Atoi(context.DefaultQuery("limit", "10"))
	if err != nil {
		context.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": "`limit` must be a valid integer"})
		return
	}
	result, err := redis.GetLeaderboardSubset(offset, limit)
	if err != nil {
		context.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"Failed to process due to error": err.Error()})
	} else {
		context.IndentedJSON(
			http.StatusOK,
			result)
	}
}
