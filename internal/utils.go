package internal

import (
	"ecampus/config"
	"encoding/base64"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/argon2"
	"log"
)

var node *snowflake.Node

// InitSnowflake initializes the Snowflake node for unique ID generation
func InitSnowflake() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Fatalf("Failed to initialize Snowflake node: %v", err)
	}
}

// GenerateID generates a unique ID using Snowflake
func GenerateID() int64 {
	return node.Generate().Int64()
}

// HashPassword hashes the given password using Argon2
func HashPassword(password string) (string, error) {
	cf, _ := config.New()
	hash := argon2.Key([]byte(password), []byte(cf.AppSecret), 3, 32*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash), nil
}

// ValidatePassword checks if the provided password matches the stored hashed password
func ValidatePassword(password, hashedPassword string) (bool, error) {
	// Decode the hashed password from base64
	hash, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, err
	}

	// Hash the provided password
	providedHash := argon2.Key([]byte(password), []byte("somesalt"), 3, 32*1024, 4, 32)

	// Compare the hashes
	if !compareHashes(hash, providedHash) {
		return false, nil
	}
	return true, nil
}

// compareHashes compares two byte slices to check if they are equal
func compareHashes(hashedPassword, providedHash []byte) bool {
	if len(hashedPassword) != len(providedHash) {
		return false
	}

	for i := range hashedPassword {
		if hashedPassword[i] != providedHash[i] {
			return false
		}
	}
	return true
}
