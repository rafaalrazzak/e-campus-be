package ecampus

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

	if match, err := helpers.VerifyPassword(creds.Password, user.Password); err != nil || !match {
		return httperror.New(http.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	userIdToken := fmt.Sprintf("%d::%s", user.ID, token)
	encryptedToken, err := helpers.EncryptAES([]byte(h.app.Config().SecretKey), []byte(userIdToken))
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	userInfo := UserInfo{
		IP:        req.RemoteAddr,
		UserAgent: req.UserAgent(),
	}

	sessionKey := bunapp.RedisKeys{}.Session(user.ID, token)
	if err := h.storeSessionData(req.Context(), sessionKey, user, token, userInfo); err != nil {
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

	userID, err := h.getUserIdFromToken(req.Context(), actualToken)
	if err != nil {
		return httperror.New(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	}

	redisKey := bunapp.RedisKeys{}.Session(userID, actualToken)
	if err := h.app.Redis().Del(req.Context(), redisKey).Err(); err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	return httpsuccess.NoContent(w, "Logged out successfully")
}

// Me retrieves the current user's session data.
func (h *AuthHandler) Me(w http.ResponseWriter, req bunrouter.Request) error {
	token := req.Header.Get("Authorization")
	if token == "" {
		return httperror.New(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	}

	actualToken := extractToken(token)
	if actualToken == "" {
		return httperror.New(http.StatusUnauthorized, "unauthorized", "Unauthorized")
	}

	decryptedToken, err := helpers.DecryptAES([]byte(h.app.Config().SecretKey), actualToken)
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	parts := strings.SplitN(string(decryptedToken), "::", 2)
	if len(parts) != 2 {
		return httperror.New(http.StatusUnauthorized, "invalid_token", "Invalid token format")
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return httperror.New(http.StatusUnauthorized, "invalid_user_id", "Invalid user ID format")
	}

	sessionKey := bunapp.RedisKeys{}.Session(userID, parts[1])

	sessionData, err := h.app.Redis().HGetAll(req.Context(), sessionKey).Result()
	if err != nil {
		return httperror.From(err, h.app.IsDebug())
	}

	return httpsuccess.OK(w, "Session data retrieved successfully", sessionData)
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

// UserInfo struct for storing user-related session information.
type UserInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

// storeSessionData stores session data in Redis.
func (h *AuthHandler) storeSessionData(ctx context.Context, sessionKey string, user *User, token string, userInfo UserInfo) error {
	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"group":      user.Group,
		"year":       user.Year,
		"token":      token,
		"role":       string(user.Role),
		"major":      string(user.Major),
		"ip":         userInfo.IP,
		"user_agent": userInfo.UserAgent,
	}

	if err := h.app.Redis().HMSet(ctx, sessionKey, sessionData).Err(); err != nil {
		return err
	}

	return h.app.Redis().Expire(ctx, sessionKey, 24*time.Hour).Err()
}

// decodeJSON decodes JSON from the request body.
func decodeJSON(body io.ReadCloser, v interface{}) error {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(body) // Close the body after decoding
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
	return fmt.Sprintf("%x", bytes), nil
}
