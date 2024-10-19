package config

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	Redis      string `env:"REDIS_URL"`
	AppSecret  string `env:"APP_SECRET"`
	Database
}

type Database struct {
	Url string `env:"DATABASE_URL"`
}
