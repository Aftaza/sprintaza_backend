package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Aftaza/sprintaza_backend/internal/model" // Sesuaikan path
	"github.com/Aftaza/sprintaza_backend/internal/utils"
)

// MockUserRepository adalah implementasi tiruan dari UserRepository.
type MockUserRepository struct {
	mock.Mock
}

// Implementasikan method dari interface UserRepository pada mock.
func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestAuthService_LoginWithPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	// Kita tidak perlu OAuth config untuk tes ini, jadi bisa nil
	authService := NewAuthService(mockRepo, nil)

	// Hash password untuk digunakan dalam tes
	hashedPassword, _ := util.HashPassword("password123")
	mockUser := &model.User{
		ID:           1,
		Email:        "user@test.com",
		PasswordHash: hashedPassword,
	}

	t.Run("Success", func(t *testing.T) {
		// Setup ekspektasi: jika FindByEmail dipanggil dengan email ini, kembalikan mockUser
		mockRepo.On("FindByEmail", "user@test.com").Return(mockUser, nil).Once()

		token, err := authService.LoginWithPassword("user@test.com", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t) // Pastikan method mock dipanggil
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Setup ekspektasi: jika FindByEmail dipanggil, kembalikan nil (tidak ditemukan)
		mockRepo.On("FindByEmail", "salah@test.com").Return(nil, nil).Once()

		token, err := authService.LoginWithPassword("salah@test.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, "kredensial tidak valid", err.Error())
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		mockRepo.On("FindByEmail", "user@test.com").Return(mockUser, nil).Once()

		token, err := authService.LoginWithPassword("user@test.com", "passwordSALAH")

		assert.Error(t, err)
		assert.Equal(t, "kredensial tidak valid", err.Error())
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}