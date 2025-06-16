package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Aftaza/sprintaza_backend/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB adalah helper untuk membuat database SQLite in-memory untuk testing.
func setupTestDB(t *testing.T) *gorm.DB {
	// "file::memory:" adalah DSN untuk SQLite in-memory
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Gagal terhubung ke database in-memory: %v", err)
	}

	// Jalankan migrasi agar tabel terbentuk
	err = db.AutoMigrate(&model.User{}, &model.UserXP{})
	if err != nil {
		t.Fatalf("Gagal melakukan migrasi: %v", err)
	}

	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	t.Run("Create User Success", func(t *testing.T) {
		// 1. Siapkan data user baru
		newUser := &model.User{
			Name:         "Budi Test",
			Email:        "budi@test.com",
			PasswordHash: "somehashedpassword",
			AuthProvider: "local",
		}

		// 2. Panggil method Create
		createdUser, err := repo.Create(newUser)

		// 3. Lakukan assertions
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.NotZero(t, createdUser.ID) // Pastikan ID sudah di-generate
		assert.Equal(t, "budi@test.com", createdUser.Email)
		
		// 4. Verifikasi bahwa UserXP juga dibuat
		assert.NotZero(t, createdUser.UserXPID)
		assert.NotNil(t, createdUser.UserXP)
		assert.Equal(t, 0, createdUser.UserXP.XP)
	})

	t.Run("Find By Email", func(t *testing.T) {
		// Kita gunakan user yang sudah dibuat di tes sebelumnya
		
		// Skenario 1: User ditemukan
		foundUser, err := repo.FindByEmail("budi@test.com")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "Budi Test", foundUser.Name)
		assert.NotNil(t, foundUser.UserXP) // Pastikan preload bekerja

		// Skenario 2: User tidak ditemukan
		notFoundUser, err := repo.FindByEmail("tidakada@test.com")
		assert.NoError(t, err) // Harusnya tidak ada error sistem, hanya record not found
		assert.Nil(t, notFoundUser)
	})
}