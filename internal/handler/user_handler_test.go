package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/crud-cicd-learn/internal/handler"
	"github.com/mohamedfawas/crud-cicd-learn/internal/service"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter is a helper function that creates a test environment
func setupTestRouter() (*gin.Engine, service.UserService) {

	// Create a new Gin router with default middleware
	router := gin.Default()

	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	// Define all the routes we want to test
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/users", userHandler.GetAllUsers)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	return router, userService
}

// TestUserHandler is our main test function that contains all user-related tests
func TestUserHandler(t *testing.T) {

	// Set Gin to test mode (reduces logging output)
	gin.SetMode(gin.TestMode)

	// Set up our test router
	router, _ := setupTestRouter()

	// Test 1: Register a new user
	t.Run("Register Success", func(t *testing.T) {
		// Create a new HTTP response recorder
		// This is like a fake HTTP response that we can inspect later
		w := httptest.NewRecorder()

		// Create the JSON request body with user registration data
		body := `{"name":"Test User","email":"test@example.com","password":"password123"}`

		// Create a new HTTP POST request
		// bytes.NewBufferString(body) converts our JSON string into a format the request can use
		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))

		// Set the content type header to indicate we're sending JSON
		req.Header.Set("Content-Type", "application/json")

		// Send our request to the router and record the response
		router.ServeHTTP(w, req)

		// Assert that we got the expected status code (201 Created)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test 2: Login with the registered user
	t.Run("Login Success", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Create login request body with the same credentials we used to register
		body := `{"email":"test@example.com","password":"password123"}`
		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		// Send the login request
		router.ServeHTTP(w, req)

		// Assert that we got a 200 OK status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Create a map to hold our response data
		var response map[string]interface{}

		// Parse the JSON response body into our map
		err := json.Unmarshal(w.Body.Bytes(), &response)

		// Assert that parsing was successful
		assert.NoError(t, err)

		// Assert that the returned user name matches what we expect
		assert.Equal(t, "Test User", response["name"])
	})

	// Test 3: Get all users
	t.Run("Get All Users", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Create a simple GET request (no body needed for GET requests)
		req, _ := http.NewRequest("GET", "/users", nil)
		router.ServeHTTP(w, req)

		// Assert we got a 200 OK status
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test 4: Get a specific user by ID
	t.Run("Get User by ID", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Create GET request for user with ID 1
		req, _ := http.NewRequest("GET", "/users/1", nil)
		router.ServeHTTP(w, req)

		// Assert we got a 200 OK status
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test 5: Update a user
	t.Run("Update User", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Create update request body with new user data
		body := `{"name":"Updated User","email":"updated@example.com"}`
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert we got a 200 OK status
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test 6: Delete a user
	t.Run("Delete User", func(t *testing.T) {
		w := httptest.NewRecorder()

		// Create DELETE request for user with ID 1
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		router.ServeHTTP(w, req)

		// Assert we got a 200 OK status
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
