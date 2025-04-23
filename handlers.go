package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

func register(context *gin.Context) {
	var newUser userRequest

	// BindJSON binds the incoming request JSON to the defined type
	if err := context.BindJSON(&newUser); err != nil {
		context.IndentedJSON(
			http.StatusBadRequest,
			err)
	}

	// Sends back registered user
	context.IndentedJSON(
		http.StatusCreated,
		newUser)
}

func addUser(username string, rdb *redis.Client, ctx *context.Context) (user, error) {
	id := uuid.NewString()
	newUser := user{
		id,
		username,
		len(username), // TODO: Fix this bc inserting rank based on username is not a good long-term option
		0,
	}

	// Implements Marshaller and can thus be converted into JSON
	serializedUser, err := json.Marshal(newUser)
	if err != nil {
		return newUser, err
	}

	err = rdb.Set(*ctx, id, serializedUser, 0).Err()
	if err != nil {
		return newUser, err
	}

	err = rdb.ZAdd(*ctx, "rank", redis.Z{
		Score:  float64(newUser.Score),
		Member: id,
	}).Err()

	if err != nil {
		return newUser, err
	}

	rank, err := rdb.ZRevRank(*ctx, "rank", id).Result()
	if err != nil {
		return newUser, err
	}
	newUser.Rank = int(rank + 1)
	return newUser, nil

}
