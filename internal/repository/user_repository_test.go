package repository

import (
	"testing"
	"errors"

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

	var createdUser *model.User

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

	t.Run("Find By ID", func(t *testing.T) {
		assert.NotNil(t, createdUser, "createdUser should not be nil")

		// Skenario 1: User ditemukan
		foundUser, err := repo.FindByID(createdUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, createdUser.Email, foundUser.Email)
		assert.NotNil(t, foundUser.UserXP) // Pastikan Preload bekerja

		// Skenario 2: User tidak ditemukan (menggunakan ID yang tidak ada)
		_, err = repo.FindByID(999)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("Update User", func(t *testing.T) {
		assert.NotNil(t, createdUser, "createdUser should not be nil")
		
		// 1. Ubah nama user yang sudah ada
		newName := "Budi Santoso"
		createdUser.Name = newName

		// 2. Panggil method Update
		err := repo.Update(createdUser)
		assert.NoError(t, err)

		// 3. Verifikasi perubahan dengan mengambil data lagi dari DB
		updatedUser, err := repo.FindByID(createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, newName, updatedUser.Name)
	})
}