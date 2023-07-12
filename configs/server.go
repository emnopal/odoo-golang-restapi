package config

import (
	"log"
	"net/http"

	env "github.com/emnopal/go_postgres/utils/env"
)

type ServerEnv struct {
	PORTS string
}

func loadServerEnv() (envConfig *ServerEnv) {
	env.LoadEnv()
	envConfig = &ServerEnv{
		PORTS: env.GetENV("PORTS", ":3000"),
	}
	return
}

func GetPort() (PORTS string) {
	envConfig := loadServerEnv()
	PORTS = envConfig.PORTS
	return
}

func ServerConfig() {
	envConfig := loadServerEnv()
	log.Fatal(http.ListenAndServe(envConfig.PORTS, nil))
}
