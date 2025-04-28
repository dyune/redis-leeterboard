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
