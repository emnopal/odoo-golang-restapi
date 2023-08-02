package config

import (
	"log"

	env "github.com/emnopal/odoo-golang-restapi/app/utils/env"
	"github.com/gin-gonic/gin"
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

func ServerConfig(r *gin.Engine) {
	envConfig := loadServerEnv()
	r.HandleMethodNotAllowed = true
	log.Fatal(r.Run(envConfig.PORTS))
}
