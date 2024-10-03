package bunapp

import "strconv"

const BaseKey = "ecampus"

// RedisKeys provides functions to generate Redis keys.
type RedisKeys struct{}

// NewRedisKeys creates a new instance of RedisKeys.
func NewRedisKeys() RedisKeys {
	return RedisKeys{}
}

// Session generates a Redis session key using user ID and session token.
func (rk RedisKeys) Session(userID int64, sessionToken string) string {
	return BaseKey + ":session:" + strconv.FormatInt(userID, 10) + ":" + sessionToken
}
