package ecampus

import (
	"errors"
	"time"

	"ecampus-be/bunapp"
)

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
	return nil
}
