package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
	DSN  string
}

// Load builds the app config from environment variables. It does not read a
// .env file itself (that's another dependency to add later if you want it,
// e.g. github.com/joho/godotenv) - for now, export the vars in your shell
// before running, or use `env $(cat .env | xargs) go run ./cmd/api`.
func Load() Config {
	port := getEnv("PORT", "8080")

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		user := getEnv("DB_USER", "root")
		pass := os.Getenv("DB_PASSWORD")
		host := getEnv("DB_HOST", "127.0.0.1")
		dbPort := getEnv("DB_PORT", "3306")
		name := getEnv("DB_NAME", "job_tracker")
		// parseTime=true is required: without it, DATE/DATETIME/TIMESTAMP
		// columns scan into []byte instead of time.Time and every scan in
		// the store layer below breaks.
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, dbPort, name)
	}

	return Config{Port: port, DSN: dsn}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
