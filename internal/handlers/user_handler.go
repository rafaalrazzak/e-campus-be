package handlers

import (
	"ecampus/database"
	"ecampus/internal/services"
	"github.com/bwmarrin/snowflake"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler is responsible for handling user-related HTTP requests
type UserHandler struct {
	service *services.UserService
	logger  *zap.Logger
}

// NewUserHandler creates a new UserHandler with the provided service
func NewUserHandler(service *services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user database.User

	if err := c.ShouldBindJSON(&user); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	// Generate a new Snowflake ID for the user
	id, _ := snowflake.NewNode(1)
	user.ID = uint64(id.Generate().Int64())

	if err := h.service.CreateUser(&user); err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetAllUsers handles GET requests to retrieve all users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser handles GET requests to retrieve a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT requests to update a user by ID
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	idStr := c.Param("id")
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}
	user.ID = uint64(int64(id))

	if err := h.service.UpdateUser(&user); err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to remove a user by ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		h.respondWithError(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// respondWithError sends a JSON response with an error message
func (h *UserHandler) respondWithError(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{"error": message})
}
