package services

import (
	"errors"

	"healthcare-app/internal/models"
	"healthcare-app/internal/repositories"
)

// Predefined errors
var (
	ErrEmailExists = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
)

// UserService handles user business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if email exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &user.ToUserResponse(), nil
}

// GetUser gets a user by ID
func (s *UserService) GetUser(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	response := user.ToUserResponse()
	return &response, nil
}

// GetAllUsers gets all users
func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToUserResponse())
	}

	return responses, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id uint, req models.CreateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if email exists for a different user
	if req.Email != user.Email {
		exists, err := s.userRepo.EmailExists(req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
	}

	// Update user
	user.Name = req.Name
	user.Email = req.Email
	user.Password = req.Password
	user.Role = req.Role

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToUserResponse()
	return &response, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(id)
} 