package ecampus

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ecampus-be/bunapp"
	"ecampus-be/ecampus/helpers"
	"ecampus-be/httputil/httperror"
	"ecampus-be/httputil/httpsuccess"
	"github.com/uptrace/bunrouter"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	app *bunapp.App
}

// NewAuthHandler creates a new AuthHandler instance.
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

func (h *AuthHandler) Login(w http.ResponseWriter, req bunrouter.Request) error {
	var creds Credentials
	if err := decodeJSON(req.Body, &creds); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	user, err := h.getUserByUsername(req.Context(), creds.Username)
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	if match, err := helpers.VerifyPassword(creds.Password, user.Password); err != nil || !match {
		return httperror.New(http.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
	}

	// Generate a random session token
	token, err := generateRandomToken(32) // 32 bytes for the token
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	// Encrypt the token
	encryptedToken, err := helpers.EncryptAES([]byte(h.app.Config().SecretKey), []byte(token))
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	sessionKey := bunapp.RedisKeys{}.Session(user.ID, token)

	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"group":      user.Group,
		"year":       user.Year,
		"role":       fmt.Sprintf("%v", user.Role),
		"major":      fmt.Sprintf("%v", user.Major),
		"ip":         req.RemoteAddr,
		"user_agent": req.UserAgent(),
	}

	// Store session data in Redis
	if err := h.app.Redis().HMSet(req.Context(), sessionKey, sessionData).Err(); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	// Set expiry for the session data
	if err := h.app.Redis().Expire(req.Context(), sessionKey, 24*time.Hour).Err(); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	return httpsuccess.Created(w, "Logged in successfully", map[string]string{"token": string(encryptedToken)})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, req bunrouter.Request) error {
	token := req.Header.Get("Authorization")
	if token == "" {
		return httperror.New(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	}

	actualToken := extractToken(token)
	if actualToken == "" {
		return httperror.New(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	}

	// Construct Redis key using the user ID and session token for deletion
	userId, err := h.getUserIdFromToken(req.Context(), actualToken)
	if err != nil {
		return httperror.New(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	}
	redisKey := bunapp.RedisKeys{}.Session(userId, actualToken)

	// Delete the session token from Redis
	if err := h.app.Redis().Del(req.Context(), redisKey).Err(); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	return httpsuccess.NoContent(w, "Logged out successfully")
}

// hashPassword hashes the user's password.
func (h *AuthHandler) hashPassword(user *User) error {
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}

// insertUser inserts a new user into the database.
func (h *AuthHandler) insertUser(ctx context.Context, user *User) error {
	_, err := h.app.DB().NewInsert().Model(user).Exec(ctx)
	return err
}

// getUserByUsername retrieves a user by their username.
func (h *AuthHandler) getUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := h.app.DB().NewSelect().Model(&user).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, err
	}
	return &user, nil
}

// decodeJSON decodes JSON from the request body.
func decodeJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close() // Close the body after decoding
	return json.NewDecoder(body).Decode(v)
}

// getUserIdFromToken retrieves the user ID associated with the session token.
func (h *AuthHandler) getUserIdFromToken(ctx context.Context, token string) (int64, error) {
	var userID int64
	if err := h.app.Redis().Get(ctx, token).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}

// extractToken extracts the actual token from the "Bearer <token>" format.
func extractToken(authHeader string) string {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}
	return ""
}

// generateRandomToken generates a random token of the specified length.
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Convert bytes to a hexadecimal string
	return fmt.Sprintf("%x", bytes), nil
}
