package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/Aftaza/sprintaza_backend/configs"
	handlerRegister "github.com/Aftaza/sprintaza_backend/handlers/auth-handlers/register"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGoogleOAuthRegister(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize database connection for testing
	db := config.Connection()
	if db == nil {
		t.Fatal("Failed to connect to database for testing")
	}

	// Create router and handler
	router := gin.New()
	registerHandler := handlerRegister.NewHandler(db)
	
	// Setup route
	router.POST("/api/v1/auth/register", registerHandler.GoogleOAuthRegister)

	t.Run("Valid Google OAuth Registration", func(t *testing.T) {
		// Create test input
		input := map[string]interface{}{
			"email": "test@example.com",
			"name":  "Test User",
		}
		
		jsonInput, _ := json.Marshal(input)
		
		// Create request
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		
		// Create response recorder
		w := httptest.NewRecorder()
		
		// Perform request
		router.ServeHTTP(w, req)
		
		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.True(t, response["success"].(bool))
		assert.Contains(t, response, "data")
		
		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "user")
		assert.Contains(t, data, "token")
		assert.True(t, data["is_new_user"].(bool))
		
		user := data["user"].(map[string]interface{})
		assert.Equal(t, "test@example.com", user["email"])
		assert.Equal(t, "Test User", user["name"])
	})

	t.Run("Invalid Input - Missing Email", func(t *testing.T) {
		// Create test input with missing email
		input := map[string]interface{}{
			"name": "Test User",
		}
		
		jsonInput, _ := json.Marshal(input)
		
		// Create request
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		
		// Create response recorder
		w := httptest.NewRecorder()
		
		// Perform request
		router.ServeHTTP(w, req)
		
		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		assert.Contains(t, response, "error")
	})

	t.Run("Invalid Input - Empty Name", func(t *testing.T) {
		// Create test input with empty name
		input := map[string]interface{}{
			"email": "test2@example.com",
			"name":  "",
		}
		
		jsonInput, _ := json.Marshal(input)
		
		// Create request
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		
		// Create response recorder
		w := httptest.NewRecorder()
		
		// Perform request
		router.ServeHTTP(w, req)
		
		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Existing User Login", func(t *testing.T) {
		// First, register a user
		input := map[string]interface{}{
			"email": "existing@example.com",
			"name":  "Existing User",
		}
		
		jsonInput, _ := json.Marshal(input)
		
		// Create first request
		req1, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonInput))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		
		// Should create user
		assert.Equal(t, http.StatusCreated, w1.Code)
		
		// Now try to register the same user again
		req2, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonInput))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		
		// Should return existing user (login)
		assert.Equal(t, http.StatusOK, w2.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w2.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		data := response["data"].(map[string]interface{})
		assert.False(t, data["is_new_user"].(bool))
		assert.Contains(t, data["message"], "already exists")
	})
}