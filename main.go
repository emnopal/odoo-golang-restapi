package main

import (
	"log"

	config "github.com/emnopal/odoo-golang-restapi/configs"
	"github.com/emnopal/odoo-golang-restapi/routes"
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
