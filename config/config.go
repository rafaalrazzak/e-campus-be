package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port             string
	DatabaseURL      string
	GoogleMapsAPIKey string
	AppSecret        string
}

func New() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Check environment variable for database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, ErrDatabaseURLNotSet
	}

	// Check environment variable for Google Maps API key
	googleMapsAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsAPIKey == "" {
		return nil, ErrGoogleMapsAPIKeyNotSet
	}

	appSecret := os.Getenv("APP_SECRET_KEY")
	if appSecret == "" {
		return nil, ErrAppSecretNotSet
	}

	return &Config{
		Port:             ":8080",
		DatabaseURL:      dbURL,
		GoogleMapsAPIKey: googleMapsAPIKey,
		AppSecret:        appSecret,
	}, nil
}

// Custom errors
var (
	ErrDatabaseURLNotSet      = errors.New("DATABASE_URL environment variable is not set")
	ErrGoogleMapsAPIKeyNotSet = errors.New("GOOGLE_MAPS_API_KEY environment variable is not set")
	ErrAppSecretNotSet        = errors.New("APP_SECRET environment variable is not set")
)
