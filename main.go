package main

import (
	"log"

	config "github.com/emnopal/go_postgres/configs"
	"github.com/emnopal/go_postgres/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	PORTS := config.GetPort()
	log.Printf("Listing for requests at http://localhost%s/", PORTS)
	r := gin.New()
	route := routes.Routes(r)
	config.ServerConfig(route)
}
