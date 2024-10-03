package ecampus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ecampus-be/bunapp"
	"ecampus-be/ecampus/helpers"
	"ecampus-be/httputil/httperror"
	"ecampus-be/httputil/httpsuccess"
	"github.com/uptrace/bunrouter"
)

type AuthHandler struct {
	app *bunapp.App
}

func NewAuthHandler(app *bunapp.App) *AuthHandler {
	return &AuthHandler{app: app}
}

// Register handles user registration.
func (h *AuthHandler) Register(w http.ResponseWriter, req bunrouter.Request) error {
	var user User
	if err := decodeJSON(req.Body, &user); err != nil {
		return httperror.BadRequest("invalid_request", "Invalid request")
	}

	user.ID = helpers.GenerateId()
	if err := h.hashPassword(&user); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	if err := h.insertUser(req.Context(), &user); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	return httpsuccess.Created(w, "User registered successfully", nil)
}

// Login handles user login.
func (h *AuthHandler) Login(w http.ResponseWriter, req bunrouter.Request) error {
	var creds Credentials
	if err := decodeJSON(req.Body, &creds); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	user, err := h.getUserByUsername(req.Context(), creds.Username)
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	if match, err := helpers.VerifyPassword(creds.Password, user.Password); err != nil {
		return httperror.From(err, h.app.IsDebug())
	} else if !match {
		return httperror.New(http.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
	}

	token, err := GenerateSessionToken(user)
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	return httpsuccess.Created(w, "Logged in successfully", map[string]string{"token": token})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, req bunrouter.Request) error {
	return httpsuccess.NoContent(w, "Logged out successfully")
}

// Helper Functions

func (h *AuthHandler) hashPassword(user *User) error {
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}

func (h *AuthHandler) insertUser(ctx context.Context, user *User) error {
	_, err := h.app.DB().NewInsert().Model(user).Exec(ctx)
	return err
}

func (h *AuthHandler) getUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := h.app.DB().NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	return &user, err
}

// decodeJSON decodes JSON from the request body.
func decodeJSON(body io.ReadCloser, v interface{}) error {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(body)
	return json.NewDecoder(body).Decode(v)
}
