package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Aftaza/sprintaza_backend/internal/config/database" // Sesuaikan path
)

func main() {
	// 1. Muat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env. Menggunakan environment variables sistem.")
	}

	// 2. Inisialisasi koneksi database
	database.ConnectDB()
	// Variabel DB global dari package database sekarang bisa diakses
	// db := database.DB 

	// 3. Inisialisasi Gin
	r := gin.Default()

	// Health Check Route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"database": "connected",
		})
	})

	// TODO: Setup Repositories, Services, Handlers (dengan mengoper db instance)
	// Contoh:
	// userRepo := repository.NewUserRepository(db)
	// userService := service.NewUserService(userRepo)
	// userHandler := handler.NewUserHandler(userService)
	// r.GET("/users", userHandler.GetUsers)

	// 4. Jalankan server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}