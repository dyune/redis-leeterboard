package redis

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var Rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "skibidi",
	DB:       0,
})

func getUserRank(id string) (int, error) {
	rank, err := Rdb.ZRevRank(ctx, "rank", id).Result()
	if err != nil {
		return -1, err
	}
	return int(rank + 1), nil
}

func AddUser(
	username string,
	score int, // ← now we pass in the desired score
) (User, error) {

	id := uuid.NewString()
	newUser := User{
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
	if err := Rdb.Set(ctx, id, data, 0).Err(); err != nil {
		return newUser, err
	}

	// Add to the sorted set "rank"
	if err := Rdb.ZAdd(ctx, "rank", redis.Z{
		Score:  float64(newUser.Score),
		Member: id,
	}).Err(); err != nil {
		return newUser, err
	}

	newUser.Rank, err = getUserRank(id)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func getUserData(id string) (User, error) {
	userData, err := Rdb.Get(ctx, id).Result()
	var existingUser User

	if err := json.Unmarshal([]byte(userData), &existingUser); err != nil {
		return User{}, err
	}

	existingUser.Rank, err = getUserRank(id)
	if err != nil {
		return User{}, err
	}
	return existingUser, nil

}
