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

func setupTestRouter() (*gin.Engine, service.UserService) {
	router := gin.Default()
	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/users", userHandler.GetAllUsers)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	return router, userService
}

func TestUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _ := setupTestRouter()

	// Test Register
	t.Run("Register Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"name":"Test User","email":"test@example.com","password":"password123"}`
		req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test Login
	t.Run("Login Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"email":"test@example.com","password":"password123"}`
		req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", response["name"])
	})

	// Test Get All Users
	t.Run("Get All Users", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test Get User by ID
	t.Run("Get User by ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test Update User
	t.Run("Update User", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := `{"name":"Updated User","email":"updated@example.com"}`
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test Delete User
	t.Run("Delete User", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
