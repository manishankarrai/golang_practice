package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("APP_ENV") // fetch this from local system // just run "export APP_ENV=dev"

	if env == "" {
		env = "dev" // default
	}

	var envFile string

	switch env {
	case "prod":
		envFile = ".env.prod"
	default:
		envFile = ".env.dev"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}
