package lib

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Println(fmt.Printf("Could not load .env file: %v", err))
	}
}

func FetchEnvVar(key string) (string, error) {
	value, ok := os.LookupEnv(key)

	if ok == false || len(value) == 0 {
		message := fmt.Sprintf("env var '%s' not found", key)
		return "", errors.New(message)
	}

	return value, nil
}
