package main

import (
	"log"

	config "github.com/emnopal/odoo-golang-restapi/app/configs"
	"github.com/emnopal/odoo-golang-restapi/app/routes"
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
