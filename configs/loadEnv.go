package config

import (
	"log"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}
