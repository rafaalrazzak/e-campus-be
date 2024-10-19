package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/matthewhartstonge/argon2"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"io"
)

func GenerateId() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node.Generate().Int64()
}

func HashData(data string) (string, error) {
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(data))
	if err != nil {
		panic(err)
	}

	return string(encoded), nil
}

func VerifyData(encodedHash, data string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(data), []byte(encodedHash))
	if err != nil {
		panic(err) // 💥
	}

	return ok, nil
}

func GenerateSessionToken() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node.Generate().Int64()
}

func GenerateSessionEncryption(sessionToken string, cfg config.Config) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(cfg.AppSecret))
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, []byte(sessionToken), nil)

	// Encode the ciphertext in base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptSessionToken(token string, cfg config.Config) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(cfg.AppSecret))
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decode the token from base64
	ciphertext, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	plaintext, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func ParseSessionToken(token string) (int64, int64, error) {
	var userId, sessionId int64
	_, err := fmt.Sscanf(token, "%d::%d", &userId, &sessionId)

	if err != nil {
		panic(err)
	}

	return userId, sessionId, nil
}
