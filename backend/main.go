package main

import (
	"fmt"
	"log"
	"net/http"

	"GopherNotes/ai"
	"GopherNotes/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Create a new Gin router instance
	r := gin.Default()

	// CORS middleware - runs before every request
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Handle preflight OPTIONS request for CORS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Otherwise, continue to the actual route handler
		c.Next()
	})

	// dummy test route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// test route for llama
	r.GET("/test-llama", func(c *gin.Context) {
		reply, err := ai.AskLlama("Say hello in one sentence")
		if err != nil {
			fmt.Println("Error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reply": reply})
	})

	// upload route
	r.POST("/upload", handlers.UploadNote)
	// chat route
	r.POST("/chat", handlers.Chat)

	r.Run(":8080")
}
