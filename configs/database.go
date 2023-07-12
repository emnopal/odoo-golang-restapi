package config

import (
	"database/sql"
	"fmt"
	"log"

	env "github.com/emnopal/go_postgres/utils/env"
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
	env.LoadEnv()
	envConfig = &DB_ENV{
		DB_CLIENT:   env.GetENV("DB_CLIENT", "postgres"),
		DB_PORTS:    env.GetENV("DB_PORTS", "5432"),
		DB_HOST:     env.GetENV("DB_HOST", "localhost"),
		DB_USER:     env.GetENV("DB_USER", "postgres"),
		DB_PASSWORD: env.GetENV("DB_PASSWORD", "postgres"),
		DB_NAME:     env.GetENV("DB_NAME", "postgres"),
	}
	return
}

func DBConfig() (*sql.DB, error) {
	dbConf := loadDBEnv()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.DB_HOST, dbConf.DB_PORTS, dbConf.DB_USER,
		dbConf.DB_PASSWORD, dbConf.DB_NAME)

	db, err := sql.Open(dbConf.DB_CLIENT, connStr)

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	return db, nil
}
