package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	Port      string
	DBPath    string
}

var AppConfig *Config

// LoadConfig loads environment variables from .env file and initializes config
func LoadConfig() {
	// Load .env file (optional - fails gracefully if not found)
	_ = godotenv.Load()

	AppConfig = &Config{
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"), // Default for development
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "rbac.db"),
	}

	log.Println("Config loaded successfully")
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
