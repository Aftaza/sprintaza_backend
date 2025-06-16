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
)

// MockAuthService adalah implementasi tiruan dari AuthService.
type MockAuthService struct {
	mock.Mock
}

// Implementasikan method dari interface AuthService pada mock.
func (m *MockAuthService) ProcessGoogleCallback(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) LoginWithPassword(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
	// Setup Gin dalam mode tes
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	// Kita tidak perlu OAuth config untuk tes ini
	handler := NewAuthHandler(mockService, nil)

	t.Run("Success", func(t *testing.T) {
		// Siapkan mock response dari service
		mockService.On("LoginWithPassword", "user@test.com", "password123").Return("fake.jwt.token", nil).Once()

		// Buat request body
		loginReq := LoginRequest{Email: "user@test.com", Password: "password123"}
		body, _ := json.Marshal(loginReq)

		// Setup recorder dan context Gin
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Panggil handler
		handler.Login(c)

		// Lakukan assertions
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "fake.jwt.token", response["token"])

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockService.On("LoginWithPassword", "user@test.com", "wrongpassword").Return("", errors.New("kredensial tidak valid")).Once()

		loginReq := LoginRequest{Email: "user@test.com", Password: "wrongpassword"}
		body, _ := json.Marshal(loginReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request Body", func(t *testing.T) {
		// Body JSON yang tidak valid
		invalidBody := []byte(`{"email":"user@test.com"}`) // password hilang

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(invalidBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// Tidak perlu mock expectation karena service tidak akan pernah dipanggil
	})
}