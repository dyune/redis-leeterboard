package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var Rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6381",
	Password: "skibidi",
	DB:       0,
})

// ZSET to hold the leaderboard: rank = { userId: score }
// SET to hold the users in the system:

func GetUserRank(id string) (int, error) {
	rank, err := Rdb.ZRevRank(ctx, "rank", id).Result()
	if err != nil {
		return -1, err
	}
	return int(rank + 1), nil
}

func AddUser(username string, score float64) (User, error) {

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

	// Store the user as a KV pair with unlimited TTL
	if err := Rdb.Set(ctx, id, data, 0).Err(); err != nil {
		return newUser, err
	}

	// Add to the sorted set "rank"
	if err := Rdb.ZAdd(ctx, "rank", redis.Z{
		Score:  newUser.Score,
		Member: id,
	}).Err(); err != nil {
		return newUser, err
	}

	newUser.Rank, err = GetUserRank(id)
	if err != nil {
		print("cock")
		return User{}, err
	}

	return newUser, nil
}

func GetUserData(id string) (User, error) {
	userData, err := Rdb.Get(ctx, id).Result()
	var existingUser User

	if err := json.Unmarshal([]byte(userData), &existingUser); err != nil {
		return User{}, err
	}

	existingUser.Rank, err = GetUserRank(id)
	if err != nil {
		return User{}, err
	}
	return existingUser, nil

}

func IncreaseUserScore(id string, scoreDelta float64) (float64, error) {

	// Increments the value of the member inside the "rank" sorted set
	score, err := Rdb.ZIncrBy(ctx, "rank", scoreDelta, id).Result()
	if err != nil {
		return -1, err
	}

	// Get the user data from Redis, unmarshal it, update it, remarshal and update it
	jsonUser, err := Rdb.Get(ctx, id).Result()
	if err != nil {
		return -1, err
	}

	var user User
	err = json.Unmarshal([]byte(jsonUser), &user)
	if err != nil {
		return -1, err
	}

	// Update the rank if need be, so fetch from the sorted set
	var rank float64
	rank, err = Rdb.ZScore(ctx, "rank", id).Result()
	if err != nil {
		return -1, err
	}

	user.Rank = int(rank)
	user.Score = score
	userData, err := json.Marshal(user)
	if err != nil {
		return -1, err
	}

	err = Rdb.Set(ctx, id, userData, 0).Err()
	if err != nil {
		return -1, err
	}

	return score, nil
}

func DeleteUser(id string) error {

	// Remove from the sorted set "rank" an entry by the ID
	err := Rdb.ZRem(ctx, "rank", id).Err()
	if err != nil {
		return err
	}

	// Then, remove the user from the overall KV pairs
	return Rdb.Del(ctx, id).Err()
}

func GetLeaderboard() ([]User, error) {

	// Get all the entries in "rank" in descending order of score
	//  0, first -> highest
	// -1, last  -> low score
	leaderboard := Rdb.ZRevRange(ctx, "rank", 0, -1).Val()
	fmt.Printf("%s", leaderboard)
	var users = make([]User, 0)

	for _, str := range leaderboard {
		var user User
		err := json.Unmarshal([]byte(str), &user)
		if err != nil {
			// For debugging purposes, print the problematic JSON string
			fmt.Printf("Error unmarshaling JSON: %v\nProblematic JSON string: %s\n", err, str)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil

}
