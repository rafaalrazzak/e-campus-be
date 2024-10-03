package ecampus

import (
	"ecampus-be/httputil/httperror"
	"ecampus-be/httputil/httpsuccess"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"ecampus-be/bunapp"
	"github.com/bwmarrin/snowflake"
	"github.com/uptrace/bunrouter"
)

type UserHandler struct {
	app *bunapp.App
}

// UserResponse is a struct for API responses, excluding sensitive fields
type UserResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Major    string `json:"major"`
	Year     int    `json:"year"`
	Phone    string `json:"phone"`
	Group    int    `json:"group"`
}

func NewUserHandler(app *bunapp.App) *UserHandler {
	return &UserHandler{app: app}
}

// List handles retrieving all users.
func (h *UserHandler) List(w http.ResponseWriter, req bunrouter.Request) error {
	var users []User
	if err := h.app.DB().NewSelect().
		Model(&users).
		Column("id", "name", "email", "phone", "username", "role", "group", "major", "year").
		Scan(req.Context()); err != nil {
		return httperror.New(http.StatusInternalServerError, "database_error", "Failed to retrieve users")
	}

	responseUsers := make([]UserResponse, len(users))
	for i, user := range users {
		responseUsers[i] = h.toUserResponse(&user)
	}

	return bunrouter.JSON(w, responseUsers)
}

// Get handles retrieving a specific user by ID.
func (h *UserHandler) Get(w http.ResponseWriter, req bunrouter.Request) error {
	id, err := h.parseID(req.Param("id"))
	if err != nil {
		return httperror.New(http.StatusBadRequest, "invalid_id", "Invalid user ID format")
	}

	var user User
	if err := h.app.DB().NewSelect().
		Model(&user).
		Column("id", "name", "email", "phone", "username", "role", "group", "major", "year").
		Where("id = ?", id).
		Scan(req.Context()); err != nil {
		return httperror.New(http.StatusNotFound, "user_not_found", "User not found")
	}

	return bunrouter.JSON(w, h.toUserResponse(&user))
}

// Update handles updating a user by ID.
func (h *UserHandler) Update(w http.ResponseWriter, req bunrouter.Request) error {
	id, err := h.parseID(req.Param("id"))
	if err != nil {
		return httperror.New(http.StatusBadRequest, "invalid_id", "Invalid user ID format")
	}

	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return httperror.New(http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
	}

	if user.ID != id {
		return httperror.New(http.StatusBadRequest, "id_mismatch", "ID in URL does not match the user ID")
	}

	if _, err := h.app.DB().NewUpdate().Model(&user).Where("id = ?", id).Exec(req.Context()); err != nil {
		return httperror.New(http.StatusInternalServerError, "update_error", "Failed to update user")
	}

	return bunrouter.JSON(w, h.toUserResponse(&user))
}

// Delete handles deleting a user by ID.
func (h *UserHandler) Delete(w http.ResponseWriter, req bunrouter.Request) error {
	id, err := h.parseID(req.Param("id"))
	if err != nil {
		return httperror.BadRequest("invalid_id", "Invalid user ID format")
	}

	if _, err := h.app.DB().NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(req.Context()); err != nil {
		return httperror.New(http.StatusInternalServerError, "delete_error", "Failed to delete user")
	}

	return httpsuccess.NoContent(w, "User deleted successfully")
}

// parseID converts a string ID to an int64.
func (h *UserHandler) parseID(idStr string) (int64, error) {
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	return id.Int64(), nil
}

// toUserResponse converts a User to UserResponse.
func (h *UserHandler) toUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		Major:    string(user.Major),
		Year:     user.Year,
		Phone:    user.Phone,
		Group:    user.Group,
	}
}

// MarshalJSON for custom JSON serialization.
func (ur UserResponse) MarshalJSON() ([]byte, error) {
	type Alias UserResponse
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    fmt.Sprintf("%d", ur.ID),
		Alias: (*Alias)(&ur),
	})
}
