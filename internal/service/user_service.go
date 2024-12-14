package service

import (
	"errors"

	"github.com/mohamedfawas/crud-cicd-learn/internal/model"
)

type UserService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest) (*model.UserResponse, error)
	GetAllUsers() ([]model.UserResponse, error)
}

type userService struct {
	users  map[string]*model.User // Simple in-memory storage for demonstration
	nextID uint
}

func NewUserService() UserService {
	return &userService{
		users:  make(map[string]*model.User),
		nextID: 1,
	}
}

func (s *userService) Register(req model.RegisterRequest) error {
	if _, exists := s.users[req.Email]; exists {
		return errors.New("user already exists")
	}

	s.users[req.Email] = &model.User{
		ID:       s.nextID,
		Email:    req.Email,
		Password: req.Password, // In real app, hash the password
	}
	s.nextID++
	return nil
}

func (s *userService) Login(req model.LoginRequest) (*model.UserResponse, error) {
	user, exists := s.users[req.Email]
	if !exists || user.Password != req.Password {
		return nil, errors.New("invalid credentials")
	}

	return &model.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (s *userService) GetAllUsers() ([]model.UserResponse, error) {
	var response []model.UserResponse
	for _, user := range s.users {
		response = append(response, model.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		})
	}
	return response, nil
}
