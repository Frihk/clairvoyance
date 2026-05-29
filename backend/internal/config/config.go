package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port    string
	AppEnv  string
	AppBaseURL string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret        string
	JWTExpiryMinutes int

	// Blockchain — all three must be set to use the real client.
	// If any is empty the router falls back to MockBlockchainClient.
	RPCUrl          string
	ContractAddress string
	PrivateKey      string

	// CORS
	AllowedOrigin string
}

var App Config

// Load reads environment variables (optionally from .env) and populates the
// package-level App configuration value.
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	expiryMinutes, err := strconv.Atoi(getEnv("JWT_EXPIRY_MINUTES", "10"))
	if err != nil {
		expiryMinutes = 10
	}

	App = Config{
		Port:        getEnv("PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		AppBaseURL:  getEnv("APP_BASE_URL", "http://localhost:8080"),

		DatabaseURL: mustGetEnv("DATABASE_URL"),

		JWTSecret:        mustGetEnv("JWT_SECRET"),
		JWTExpiryMinutes: expiryMinutes,

		RPCUrl:          getEnv("RPC_URL", ""),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		PrivateKey:      getEnv("PRIVATE_KEY", ""),

		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}