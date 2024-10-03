package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Parameters for Argon2 hashing
type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// HashPassword hashes the password using Argon2 and returns the hash.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16) // 16-byte salt

	// Generate a random salt
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	p := params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	// Hash the password using Argon2
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Encode salt and hash for storage
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash))

	return encodedHash, nil
}

// VerifyPassword checks if the given password matches the stored hash.
func VerifyPassword(password, encodedHash string) (bool, error) {
	p, salt, originalHash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash the provided password using the same parameters
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Compare the original hash with the new hash using constant-time comparison
	if subtle.ConstantTimeCompare(originalHash, hash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHash extracts the parameters, salt, and hash from the encoded hash.
func decodeHash(encodedHash string) (*params, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil || version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible Argon2 version")
	}

	p := &params{}
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}

	p.saltLength = uint32(len(salt))
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// EnsureKeySize ensures that the AES key is of the correct length (16, 24, or 32 bytes).
func EnsureKeySize(key []byte) ([]byte, error) {
	switch len(key) {
	case 16, 24, 32:
		return key, nil
	default:
		return nil, errors.New("invalid AES key size (must be 16, 24, or 32 bytes)")
	}
}

// EncryptAES encrypts plaintext using AES-GCM and returns a base64 encoded ciphertext.
func EncryptAES(key, plaintext []byte) (string, error) {
	// Ensure the key size is valid for AES
	key, err := EnsureKeySize(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data using AES-GCM
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts a base64 encoded ciphertext using AES-GCM and returns the plaintext.
func DecryptAES(key []byte, encryptedData string) ([]byte, error) {
	// Ensure the key size is valid for AES
	key, err := EnsureKeySize(key)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
