package config

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"3000"`
	Database
}

type Database struct {
	Url string `env:"DATABASE_URL"`
}
