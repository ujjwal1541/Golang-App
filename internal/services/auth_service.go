package services

import (
	"errors"
	"fmt"
	"time"

	"healthcare-app/internal/models"
	"healthcare-app/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
)

// Predefined errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenGeneration    = errors.New("could not generate token")
	ErrInvalidToken       = errors.New("invalid token")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Claims represents the JWT claims
type Claims struct {
	UserID uint           `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// Login authenticates a user
func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !models.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	return &models.LoginResponse{
		Token: token,
		User:  user.ToUserResponse(),
	}, nil
}

// GenerateToken generates a JWT token
func (s *AuthService) GenerateToken(userID uint, role models.UserRole) (string, error) {
	// Create claims with expiration
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
} 