package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/crud-cicd-learn/internal/handler"
	"github.com/mohamedfawas/crud-cicd-learn/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(req model.RegisterRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockUserService) Login(req model.LoginRequest) (*model.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) GetAllUsers() ([]model.UserResponse, error) {
	args := m.Called()
	return args.Get(0).([]model.UserResponse), args.Error(1)
}

func TestRegister(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	tests := []struct {
		name         string
		input        model.RegisterRequest
		setupMock    func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success",
			input: model.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("Register", mock.AnythingOfType("model.RegisterRequest")).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"User registered successfully"}`,
		},
		{
			name: "User Already Exists",
			input: model.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("Register", mock.AnythingOfType("model.RegisterRequest")).
					Return(errors.New("user already exists"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"user already exists"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockService.ExpectedCalls = nil
			tt.setupMock()

			// Create request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonValue, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
			c.Request.Header.Set("Content-Type", "application/json")

			// Make request
			userHandler.Register(c)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	successResponse := &model.UserResponse{
		ID:    1,
		Email: "test@example.com",
	}

	tests := []struct {
		name         string
		input        model.LoginRequest
		setupMock    func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success",
			input: model.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("Login", mock.AnythingOfType("model.LoginRequest")).
					Return(successResponse, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"email":"test@example.com"}`,
		},
		{
			name: "Invalid Credentials",
			input: model.LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpass",
			},
			setupMock: func() {
				mockService.On("Login", mock.AnythingOfType("model.LoginRequest")).
					Return(nil, errors.New("invalid credentials"))
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"invalid credentials"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonValue, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
			c.Request.Header.Set("Content-Type", "application/json")

			userHandler.Login(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
