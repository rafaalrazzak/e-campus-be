package ecampus

import (
	"encoding/json"
	"errors"
	"fmt"
	"ecampus-be/ecampus/helpers"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"ecampus-be/bunapp"
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

func (h *UserHandler) List(w http.ResponseWriter, req bunrouter.Request) error {
	var users []User
	if err := h.app.DB().NewSelect().
		Model(&users).
		Column("id", "name", "email", "phone", "username", "role", "group", "major", "year").
		Scan(req.Context()); err != nil {
		return h.handleError(w, err, http.StatusInternalServerError)
	}

	responseUsers := make([]UserResponse, len(users))
	for i, user := range users {
		responseUsers[i] = h.toUserResponse(&user)
	}

	return bunrouter.JSON(w, responseUsers)
}

func (h *UserHandler) Get(w http.ResponseWriter, req bunrouter.Request) error {
	id, err := h.parseID(req.Param("id"))
	if err != nil {
		return h.handleError(w, err, http.StatusBadRequest)
	}

	var user User
	if err := h.app.DB().NewSelect().
		Model(&user).
		Column("id", "name", "email", "phone", "username", "role", "group", "major", "year").
		Where("id = ?", id).
		Scan(req.Context()); err != nil {
		return h.handleError(w, err, http.StatusNotFound)
	}

	return bunrouter.JSON(w, h.toUserResponse(&user))
}

func (h *UserHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return h.handleError(w, err, http.StatusBadRequest)
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return h.handleError(w, err, http.StatusInternalServerError)
	}
	user.ID = node.Generate().Int64()
	user.Password, err = helpers.PasswordHasher
	if err != nil {
		return h.handleError(w, err, http.StatusInternalServerError)
	}

	if _, err := h.app.DB().NewInsert().Model(&user).Exec(req.Context()); err != nil {
		return h.handleError(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, h.toUserResponse(&user))
}

func (h *UserHandler) Update(w http.ResponseWriter, req bunrouter.Request) error {
	id, err := h.parseID(req.Param("id"))
	if err != nil {
		return h.handleError(w, err, http.StatusBadRequest)
	}

	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		return h.handleError(w, err, http.StatusBadRequest)
	}

	if user.ID != id {
		return h.handleError(w, errors.New("ID mismatch"), http.StatusBadRequest)
	}

	if _, err := h.app.DB().NewUpdate().Model(&user).Where("id = ?", id).Exec(req.Context()); err != nil {
		return h.handleError(w, err, http.StatusInternalServerError)
	}

	return bunrouter.JSON(w, h.toUserResponse(&user))
}

func (h *UserHandler) Delete(w http.ResponseWriter, req bunrouter.Request) error {
	id, err := h.parseID(req.Param("id"))
	if err != nil {
		return h.handleError(w, err, http.StatusBadRequest)
	}

	if _, err := h.app.DB().NewDelete().Model((*User)(nil)).Where("id = ?", id).Exec(req.Context()); err != nil {
		return h.handleError(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *UserHandler) parseID(idStr string) (int64, error) {
	id, err := snowflake.ParseString(idStr)
	if err != nil {
		return 0, errors.New("invalid user ID format")
	}
	return id.Int64(), nil
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error, status int) error {
	http.Error(w, err.Error(), status)
	return err
}

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
