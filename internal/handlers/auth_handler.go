package handlers

import (
	"errors"
	"net/http"
	"strings"

	"healthcare-app/internal/models"
	"healthcare-app/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles login requests
// @Summary Login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// RequireAuth is a middleware to require authentication
func (h *AuthHandler) RequireAuth(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Missing authorization header"})
			c.Abort()
			return
		}

		// Check if the header is in the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := h.authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		// Set the user in the context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		next(c)
	}
}

// RequireReceptionist is a middleware to require receptionist role
func (h *AuthHandler) RequireReceptionist(c *gin.Context) {
	role := c.GetString("userRole")
	if role != string(models.RoleReceptionist) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Receptionist role required"})
		c.Abort()
		return
	}
	c.Next()
}

// RequireDoctor is a middleware to require doctor role
func (h *AuthHandler) RequireDoctor(c *gin.Context) {
	role := c.GetString("userRole")
	if role != string(models.RoleDoctor) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Doctor role required"})
		c.Abort()
		return
	}
	c.Next()
} 