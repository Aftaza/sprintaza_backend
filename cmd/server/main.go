package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/Aftaza/sprintaza_backend/internal/config/database" 
	"github.com/Aftaza/sprintaza_backend/internal/handler"
)

func main() {
	// 1. Muat variabel dari file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env.")
	}

	// 2. Inisialisasi koneksi database (termasuk migrasi dan seeding)
	database.ConnectDB()

	// 3. Setup semua route dan dependencies
	r := handler.SetupRoutes(database.DB)

	// 4. Jalankan server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Tidak dapat menjalankan server: %v", err)
	}
}