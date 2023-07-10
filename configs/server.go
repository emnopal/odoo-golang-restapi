package config

import (
	"log"
	"net/http"

	helper "github.com/emnopal/go_helper"
)

type ServerEnv struct {
	PORTS string
}

func loadServerEnv() (envConfig *ServerEnv) {
	loadEnv()

	envConfig = &ServerEnv{
		PORTS: helper.GetENV("PORTS", ":3000"),
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
