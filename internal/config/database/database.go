package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Aftaza/sprintaza_backend/internal/model" 
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB menginisialisasi koneksi, migrasi, dan seeding database.
func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Koneksi database berhasil.")

	// Auto-migrate semua model
	err = DB.AutoMigrate(
		&model.User{}, &model.UserXP{}, &model.Project{}, &model.Role{},
		&model.ProjectMember{}, &model.Column{}, &model.Task{},
		&model.Achievement{}, &model.UserAchievement{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	log.Println("Migrasi database berhasil.")

	// Panggil fungsi seeding setelah migrasi
	seedDatabase(DB)
}

// seedDatabase adalah fungsi utama untuk memanggil semua seeder.
func seedDatabase(db *gorm.DB) {
	log.Println("Memulai proses seeding...")
	seedRoles(db)
	seedAchievements(db)
	log.Println("Proses seeding selesai.")
}

// seedRoles mengisi tabel 'roles' jika kosong.
func seedRoles(db *gorm.DB) {
	// 1. Cek apakah tabel roles sudah ada isinya
	var count int64
	db.Model(&model.Role{}).Count(&count)
	if count > 0 {
		log.Println("Tabel 'roles' sudah berisi data, seeding dilewati.")
		return
	}

	// 2. Jika kosong, siapkan data seed
	roles := []model.Role{
		{Name: "Owner", Description: "Memiliki akses penuh terhadap proyek, termasuk pengaturan dan penghapusan."},
		{Name: "Admin", Description: "Memiliki akses fungsional terhadap proyek (create, read, dan update), dibawah tingkat owner."},
		{Name: "Member", Description: "Anggota tim yang dapat mengerjakan tugas di dalam proyek."},
	}

	// 3. Masukkan data ke database
	if err := db.Create(&roles).Error; err != nil {
		log.Fatalf("Gagal melakukan seeding untuk tabel 'roles': %v", err)
	}
	log.Println("Seeding tabel 'roles' berhasil.")
}

// seedAchievements mengisi tabel 'achievements' jika kosong.
func seedAchievements(db *gorm.DB) {
	// 1. Cek apakah tabel achievements sudah ada isinya
	var count int64
	db.Model(&model.Achievement{}).Count(&count)
	if count > 0 {
		log.Println("Tabel 'achievements' sudah berisi data, seeding dilewati.")
		return
	}

	// 2. Jika kosong, siapkan data seed berdasarkan ERD
	achievements := []model.Achievement{
		{Name: "Selamat Datang!", Description: "Berhasil mendaftar dan masuk untuk pertama kalinya.", XpReward: 10, IconURL: "/icons/welcome.png"},
		{Name: "Perencana Proyek", Description: "Berhasil membuat proyek pertama Anda.", XpReward: 25, IconURL: "/icons/project_planner.png"},
		{Name: "Tugas Pertama", Description: "Berhasil menyelesaikan tugas pertama Anda.", XpReward: 15, IconURL: "/icons/first_task.png"},
		{Name: "Anggota Tim", Description: "Berhasil bergabung ke dalam sebuah proyek.", XpReward: 5, IconURL: "/icons/team_member.png"},
		{Name: "Papan Baru", Description: "Berhasil membuat kolom (board) baru di dalam proyek.", XpReward: 5, IconURL: "/icons/new_board.png"},
		{Name: "Koneksi Kode", Description: "Berhasil mengintegrasikan repositori Git ke proyek.", XpReward: 50, IconURL: "/icons/code_connect.png"},
		{Name: "Delegator Andal", Description: "Berhasil mendelegasikan 10 tugas kepada anggota tim lain.", XpReward: 20, IconURL: "/icons/delegator.png"},
		{Name: "Got Your Back!", Description: "Berhasil membantu anggota lain menyelesaikan tugas yang dialihkan.", XpReward: 35, IconURL: "/icons/got_your_back.png"},
		{Name: "Kontributor", Description: "Berhasil menyelesaikan 25 tugas.", XpReward: 80, IconURL: "/icons/contributor.png"},
	}

	// 3. Masukkan data ke database
	if err := db.Create(&achievements).Error; err != nil {
		log.Fatalf("Gagal melakukan seeding untuk tabel 'achievements': %v", err)
	}

	log.Println("Seeding tabel 'achievements' berhasil.")
}