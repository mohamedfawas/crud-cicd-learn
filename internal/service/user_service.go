package service

import (
	"errors"

	"github.com/mohamedfawas/crud-cicd-learn/internal/model"
)

type UserService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest) (*model.UserResponse, error)
	GetAllUsers() ([]model.UserResponse, error)
	GetUserByID(id uint) (*model.UserResponse, error)
	UpdateUser(id uint, req model.UpdateUserRequest) error
	DeleteUser(id uint) error
}

type MockUserService struct {
	users  map[uint]*model.User
	nextID uint
}

func NewUserService() UserService {
	return &MockUserService{
		users:  make(map[uint]*model.User),
		nextID: 1,
	}
}

func (s *MockUserService) Register(req model.RegisterRequest) error {
	// Check if email already exists
	for _, user := range s.users {
		if user.Email == req.Email {
			return errors.New("email already exists")
		}
	}

	// Create new user
	user := &model.User{
		ID:       s.nextID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	s.users[s.nextID] = user
	s.nextID++
	return nil
}

func (s *MockUserService) Login(req model.LoginRequest) (*model.UserResponse, error) {
	for _, user := range s.users {
		if user.Email == req.Email && user.Password == req.Password {
			return &model.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			}, nil
		}
	}
	return nil, errors.New("invalid credentials")
}

func (s *MockUserService) GetAllUsers() ([]model.UserResponse, error) {
	var response []model.UserResponse
	for _, user := range s.users {
		response = append(response, model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return response, nil
}

func (s *MockUserService) GetUserByID(id uint) (*model.UserResponse, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *MockUserService) UpdateUser(id uint, req model.UpdateUserRequest) error {
	user, exists := s.users[id]
	if !exists {
		return errors.New("user not found")
	}

	// Check if new email already exists
	for _, u := range s.users {
		if u.ID != id && u.Email == req.Email {
			return errors.New("email already exists")
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	return nil
}

func (s *MockUserService) DeleteUser(id uint) error {
	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(s.users, id) // deletes key-value pair from the map
	return nil
}
