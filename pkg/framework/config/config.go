package config

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	Redis      string `env:"REDIS_URL"`
	Database
}

type Database struct {
	Url string `env:"DATABASE_URL"`
}
