package config

import (
	"database/sql"
	"fmt"
	"log"

	helper "github.com/emnopal/go_helper"
	_ "github.com/lib/pq"
)

type DB_ENV struct {
	DB_CLIENT   string
	DB_PORTS    string
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

func loadDBEnv() (envConfig *DB_ENV) {

	loadEnv()

	envConfig = &DB_ENV{
		DB_CLIENT:   helper.GetENV("DB_CLIENT", "postgres"),
		DB_PORTS:    helper.GetENV("DB_PORTS", "5432"),
		DB_HOST:     helper.GetENV("DB_HOST", "localhost"),
		DB_USER:     helper.GetENV("DB_USER", "postgres"),
		DB_PASSWORD: helper.GetENV("DB_PASSWORD", "postgres"),
		DB_NAME:     helper.GetENV("DB_NAME", "postgres"),
	}
	return
}

func DBConfig() (db *sql.DB, err error) {

	envConfig := loadDBEnv()

	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		envConfig.DB_CLIENT, envConfig.DB_USER,
		envConfig.DB_PASSWORD, envConfig.DB_PORTS,
		envConfig.DB_HOST, envConfig.DB_NAME,
	)

	fmt.Println(connStr)

	db, err = sql.Open(envConfig.DB_CLIENT, connStr)

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	return
}
