package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/matthewhartstonge/argon2"
)

func GenerateId() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	id := node.Generate()

	return id.Int64()
}
func HashData(data string) string {
	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(data))
	if err != nil {
		panic(err)
	}

	return string(hash)
}

func VerifyData(hash string, data string) bool {
	match, err := argon2.VerifyEncoded([]byte(data), []byte(hash))
	if err != nil {
		panic(err)
	}

	return match
}

func GenerateSessionToken() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}

	id := node.Generate()

	return id.String(), nil
}
