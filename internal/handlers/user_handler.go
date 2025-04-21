package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"healthcare-app/internal/models"
	"healthcare-app/internal/services"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser handles create user requests
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Create User Request"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrEmailExists) {
			status = http.StatusBadRequest
		}
		RespondWithError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles get user requests
// @Summary Get user
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrUserNotFound) {
			status = http.StatusNotFound
		}
		RespondWithError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers handles get all users requests
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {array} models.UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser handles update user requests
// @Summary Update user
// @Description Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.CreateUserRequest true "Update User Request"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrUserNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, services.ErrEmailExists) {
			status = http.StatusBadRequest
		}
		RespondWithError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles delete user requests
// @Summary Delete user
// @Description Delete a user
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrUserNotFound) {
			status = http.StatusNotFound
		}
		RespondWithError(c, status, err.Error())
		return
	}

	RespondWithSuccess(c, "User deleted successfully", nil)
} 