package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	// Init router
	router := gin.Default()
	fmt.Println("Running leaderboard...")

	// Uncomment for release: gin.SetMode(gin.ReleaseMode)

	// Init endpoints and run
	api := router.Group("/api/v1")
	setupRoutes(api)
	err := router.Run(":8080")

	if err != nil {
		log.Fatal("Application failure: ", err)
	}

}
