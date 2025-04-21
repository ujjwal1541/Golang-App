package services

import (
	"errors"
	"testing"

	"healthcare-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)
	
	// Create test user with hashed password
	hashedPassword, _ := models.HashPassword("password123")
	user := &models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     models.RoleReceptionist,
	}
	
	// Setup expectations
	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
	
	// Create service with mock
	service := NewAuthService(mockRepo, "test-secret")
	
	// Test login
	req := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	res, err := service.Login(req)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, user.ID, res.User.ID)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)
	
	// Create test user with hashed password
	hashedPassword, _ := models.HashPassword("password123")
	user := &models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     models.RoleReceptionist,
	}
	
	// Setup expectations
	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
	
	// Create service with mock
	service := NewAuthService(mockRepo, "test-secret")
	
	// Test login with wrong password
	req := models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong-password",
	}
	
	res, err := service.Login(req)
	
	// Assert results
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, ErrInvalidCredentials, err)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)
	
	// Setup expectations
	mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))
	
	// Create service with mock
	service := NewAuthService(mockRepo, "test-secret")
	
	// Test login with non-existent user
	req := models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	
	res, err := service.Login(req)
	
	// Assert results
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, ErrInvalidCredentials, err)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestValidateToken(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)
	
	// Create service with mock
	service := NewAuthService(mockRepo, "test-secret")
	
	// Generate a token
	token, err := service.GenerateToken(1, models.RoleReceptionist)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Validate the token
	claims, err := service.ValidateToken(token)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, models.RoleReceptionist, claims.Role)
} 