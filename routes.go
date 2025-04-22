package main

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.RouterGroup) {

	// GET: Information on application
	router.GET("/", about)

	// POST: Register a user
	router.POST("/users", register)
}
