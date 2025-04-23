package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "skibidi",
	DB:       0,
})

func getUserRank(id string) (int, error) {
	rank, err := rdb.ZRevRank(ctx, "rank", id).Result()
	if err != nil {
		return -1, err
	}
	return int(rank + 1), nil
}

func addUser(
	username string,
	score int, // ← now we pass in the desired score
) (user, error) {

	id := uuid.NewString()
	newUser := user{
		Id:    id,
		Name:  username,
		Score: score, // ← use the passed-in score
		Rank:  0,     // will fill in below
	}

	// Marshal to JSON
	data, err := json.Marshal(newUser)
	if err != nil {
		return newUser, err
	}

	// Store the user object with no TTL and no expiration
	if err := rdb.Set(ctx, id, data, 0).Err(); err != nil {
		return newUser, err
	}

	// Add to the sorted set "rank"
	if err := rdb.ZAdd(ctx, "rank", redis.Z{
		Score:  float64(newUser.Score),
		Member: id,
	}).Err(); err != nil {
		return newUser, err
	}

	newUser.Rank, err = getUserRank(id)
	if err != nil {
		return user{}, err
	}

	return newUser, nil
}

func getUserData(id string) (user, error) {
	userData, err := rdb.Get(ctx, id).Result()
	var existingUser user

	if err := json.Unmarshal([]byte(userData), &existingUser); err != nil {
		return user{}, err
	}

	existingUser.Rank, err = getUserRank(id)
	if err != nil {
		return user{}, err
	}
	return existingUser, nil

}
