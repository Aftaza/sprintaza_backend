package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Aftaza/sprintaza_backend/internal/model"
	"github.com/Aftaza/sprintaza_backend/internal/utils"
)

// MockUserRepository dari auth_service_test.go bisa digunakan lagi,
// kita hanya perlu menambahkan implementasi method baru.

// type MockUserRepository struct {
// 	mock.Mock
// }

// func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) { /* ... */ return nil, nil }
// func (m *MockUserRepository) Create(user *model.User) (*model.User, error) { /* ... */ return nil, nil }

// Implementasi method baru untuk mock
func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}


func TestUserService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	hashedPassword, _ := util.HashPassword("passwordLama123")
	mockUser := &model.User{ID: 1, Name: "User Lama", Email: "user@test.com", PasswordHash: hashedPassword}

	t.Run("GetProfile Success", func(t *testing.T) {
		mockRepo.On("FindByID", uint(1)).Return(mockUser, nil).Once()

		profile, err := userService.GetProfile(1)
		
		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, mockUser.Name, profile.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateName Success", func(t *testing.T) {
		// Ekspektasi: FindByID akan dipanggil, lalu Update akan dipanggil
		mockRepo.On("FindByID", uint(1)).Return(mockUser, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil).Once()

		updatedUser, err := userService.UpdateName(1, "User Baru")

		assert.NoError(t, err)
		assert.Equal(t, "User Baru", updatedUser.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePassword Success", func(t *testing.T) {
		mockRepo.On("FindByID", uint(1)).Return(mockUser, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil).Once()

		err := userService.UpdatePassword(1, "passwordLama123", "passwordBaru456")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePassword Failed - Wrong Old Password", func(t *testing.T) {
		// Hanya FindByID yang akan dipanggil, Update tidak akan pernah tercapai.
		mockRepo.On("FindByID", uint(1)).Return(mockUser, nil).Once()

		err := userService.UpdatePassword(1, "passwordLamaSALAH", "passwordBaru456")

		assert.Error(t, err)
		assert.Equal(t, "password lama tidak sesuai", err.Error())
		// Pastikan hanya method yang diharapkan yang dipanggil
		mockRepo.AssertExpectations(t)
	})
}