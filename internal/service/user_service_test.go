package service_test

import (
	"testing"

	"github.com/mohamedfawas/crud-cicd-learn/internal/model"
	"github.com/mohamedfawas/crud-cicd-learn/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	svc := service.NewUserService()

	// Test Register
	t.Run("Register Success", func(t *testing.T) {
		req := model.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		err := svc.Register(req)
		assert.NoError(t, err)
	})

	t.Run("Register Duplicate Email", func(t *testing.T) {
		req := model.RegisterRequest{
			Name:     "Another User",
			Email:    "test@example.com", // Same email
			Password: "password456",
		}
		err := svc.Register(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})

	// Test Login
	t.Run("Login Success", func(t *testing.T) {
		req := model.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		user, err := svc.Login(req)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("Login Invalid Credentials", func(t *testing.T) {
		req := model.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		_, err := svc.Login(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	// Test Get User
	t.Run("Get User Success", func(t *testing.T) {
		user, err := svc.GetUserByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", user.Name)
	})

	t.Run("Get User Not Found", func(t *testing.T) {
		_, err := svc.GetUserByID(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	// Test Update User
	t.Run("Update User Success", func(t *testing.T) {
		req := model.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		err := svc.UpdateUser(1, req)
		assert.NoError(t, err)

		user, _ := svc.GetUserByID(1)
		assert.Equal(t, "Updated Name", user.Name)
		assert.Equal(t, "updated@example.com", user.Email)
	})

	t.Run("Update User Not Found", func(t *testing.T) {
		req := model.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		err := svc.UpdateUser(999, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Update User Duplicate Email", func(t *testing.T) {
		// First register another user
		err := svc.Register(model.RegisterRequest{
			Name:     "Another User",
			Email:    "another@example.com",
			Password: "password123",
		})
		assert.NoError(t, err)

		// Try to update first user with second user's email
		req := model.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "another@example.com",
		}
		err = svc.UpdateUser(1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})

	// Test Delete User
	t.Run("Delete User Success", func(t *testing.T) {
		err := svc.DeleteUser(1)
		assert.NoError(t, err)

		_, err = svc.GetUserByID(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Delete User Not Found", func(t *testing.T) {
		err := svc.DeleteUser(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}
