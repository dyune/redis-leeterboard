package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"redis-leetcode-leaderboard/api"
)

func main() {

	// Init router
	router := gin.Default()
	fmt.Println("Running leaderboard...")

	// Uncomment for release: gin.SetMode(gin.ReleaseMode)

	// Init endpoints and run
	apiRoutes := router.Group("/api/v1")
	api.SetupRoutes(apiRoutes)
	err := router.Run(":8080")

	if err != nil {
		log.Fatal("Application failure: ", err)
	}

}
