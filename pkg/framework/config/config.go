package config

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	Database
}

type Database struct {
	Url string `env:"DATABASE_URL"`
}
