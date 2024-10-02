package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
)

// HashPassword hashes the password using Argon2.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16) // Generate a random 16-byte salt
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password using Argon2
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Combine the salt and hash for storage
	hashWithSalt := append(salt, hash...)
	return base64.RawStdEncoding.EncodeToString(hashWithSalt), nil
}

// VerifyPassword verifies the hashed password with the plaintext password.
func VerifyPassword(hashedPassword, password string) (bool, error) {
	// Decode the stored password hash
	data, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, err
	}

	if len(data) < 16+32 { // 16 bytes for salt + 32 bytes for hash
		return false, fmt.Errorf("hashed password is too short")
	}

	// Extract the salt and the original hash
	salt := data[:16]
	originalHash := data[16:]

	// Hash the provided password using the same salt
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Compare the original hash with the new hash
	return string(originalHash) == string(hash), nil
}
