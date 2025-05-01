package tests

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	redisModule "redis-leetcode-leaderboard/redis"
	"testing"
)

var testCtx = context.Background()
var testRdb *redis.Client

// TestMain handles setup and teardown for all tests
func TestMain(m *testing.M) {
	// Setup test Redis client
	testRdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "skibidi", // Use your Redis password
		DB:       1,         // Use a different DB than production
	})

	// Replace the global Redis client with our test client
	redisModule.Rdb = testRdb

	// Run tests
	code := m.Run()

	// Clean up after tests
	testRdb.FlushDB(testCtx)
	err := testRdb.Close()
	if err != nil {
		return
	}

	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	// Clean DB before test
	testRdb.FlushDB(testCtx)

	// Test data
	username := "testuser"
	score := 100.00

	// Call function to test
	newUser, err := redisModule.AddUser(username, score)

	// Assertions
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}

	if newUser.Name != username {
		t.Errorf("Expected username %s, got %s", username, newUser.Name)
	}

	if newUser.Score != score {
		t.Errorf("Expected score %f, got %f", score, newUser.Score)
	}

	if newUser.Rank != 1 {
		t.Errorf("Expected rank 1, got %d", newUser.Rank)
	}

	// Verify data was stored in Redis
	exists, err := testRdb.Exists(testCtx, newUser.Id).Result()
	if err != nil || exists != 1 {
		t.Errorf("User data not stored in Redis")
	}

	// Verify user was added to the leaderboard
	score, err = testRdb.ZScore(testCtx, "rank", newUser.Id).Result()
	if err != nil || score != newUser.Score {
		t.Errorf("User not properly added to leaderboard")
	}
}
