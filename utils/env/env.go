package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetENV(key string, default_value string) (ENV_VALUE string) {
	ENV_VALUE = os.Getenv(key)
	if ENV_VALUE == "" {
		ENV_VALUE = default_value
	}
	return
}

func LoadEnv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
