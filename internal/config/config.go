package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
	JWT  JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type HTTPConfig struct {
	Port string
}

type JWTConfig struct {
	SecretKey string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Если .env файл не найден, продолжаем с переменными окружения
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "forum"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", ":8080"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
