package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Aftaza/sprintaza_backend/internal/model"
)

// MockUserService adalah implementasi tiruan dari UserService.
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetProfile(userID uint) (*model.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserService) UpdateName(userID uint, newName string) (*model.User, error) {
	args := m.Called(userID, newName)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}


func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockUser := &model.User{ID: 1, Name: "Test User", Email: "test@user.com"}

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetProfile", uint(1)).Return(mockUser, nil).Once()
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Simulasikan middleware yang men-set userID
		c.Set("userID", uint(1))

		handler.GetProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response model.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, mockUser.Name, response.Name)

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		mockService.On("UpdatePassword", uint(1), "oldPass", "newPass").Return(nil).Once()

		reqBody := UpdatePasswordRequest{OldPassword: "oldPass", NewPassword: "newPass"}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/me/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", uint(1))

		handler.UpdatePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Failed - Incorrect Old Password", func(t *testing.T) {
		mockService.On("UpdatePassword", uint(1), "wrongOldPass", "newPass").Return(errors.New("password lama tidak sesuai")).Once()

		reqBody := UpdatePasswordRequest{OldPassword: "wrongOldPass", NewPassword: "newPass"}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/me/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", uint(1))

		handler.UpdatePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}