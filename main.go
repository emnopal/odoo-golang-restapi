package cmd

import (
	"log"

	config "github.com/emnopal/go_postgres/configs"
	"github.com/emnopal/go_postgres/routes"
	_ "github.com/lib/pq"
)

func RunMainService() {
	PORTS := config.GetPort()
	log.Printf("Listing for requests at http://localhost%s/", PORTS)
	routes.Routes()
	config.ServerConfig()
}
