package ecampus

import (
	"encoding/json"
	"errors"
	"time"

	"ecampus-be/bunapp"
	"ecampus-be/ecampus/helpers"
)

var encryptionKey []byte

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
	Role     Role      `json:"role"`
	ExpireAt time.Time `json:"expire_at"`
}

// InitAuth initializes the authentication system with the provided configuration.
func InitAuth(cfg *bunapp.AppConfig) error {
	if cfg.SecretKey == "" {
		return errors.New("secret key is not set in the configuration")
	}
	encryptionKey = []byte(cfg.SecretKey)
	return nil
}

// GenerateSessionToken creates a new session token for the user.
func GenerateSessionToken(user *User) (string, error) {
	session := Session{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	return helpers.EncryptAES(encryptionKey, sessionJSON)
}

// DecryptSessionToken decrypts and validates the session token.
func DecryptSessionToken(tokenString string) (*Session, error) {
	decrypted, err := helpers.DecryptAES(encryptionKey, tokenString)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(decrypted, &session); err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpireAt) {
		return nil, errors.New("session expired")
	}

	return &session, nil
}
