package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to resolve config file path")
	}

	envPath := filepath.Join(filepath.Dir(filename), "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("No .env file found at %s, using existing environment variables", envPath)
	}
}
