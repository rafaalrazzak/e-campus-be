package bunapp

import (
	"embed"
	"io/fs"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	//go:embed embed
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

// FS returns the embedded file system.
func FS() fs.FS {
	unwrapFSOnce.Do(func() {
		fsys, err := fs.Sub(embedFS, "embed")
		if err != nil {
			panic(err)
		}
		unwrappedFS = fsys
	})
	return unwrappedFS
}

// AppConfig holds the application configuration.
type AppConfig struct {
	Service   string
	Env       string
	Debug     bool
	SecretKey string
	DbUrl     string
	RedisUrl  string
}

// LoadConfig loads the configuration from a .env file.
func LoadConfig(service, env string) (*AppConfig, error) {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &AppConfig{
		Service:   service,
		Env:       env,
		Debug:     os.Getenv("DEBUG") == "true",
		SecretKey: os.Getenv("SECRET_KEY"),
		DbUrl:     os.Getenv("DATABASE_URL"),
		RedisUrl:  os.Getenv("REDIS_URL"),
	}

	return cfg, nil
}
