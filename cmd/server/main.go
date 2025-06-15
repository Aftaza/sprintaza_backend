package main

import (
	"log"

	"github.com/gin-gonic/gin"
	// kelak import config, router, dll dari folder internal
)

func main() {
	// Inisialisasi Gin
	r := gin.Default()

	// Health Check Route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// TODO: Load Config
	// TODO: Setup Database Connection
	// TODO: Setup Repositories, Services, Handlers
	// TODO: Setup Routes from internal/handler/router.go

	// Jalankan server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}