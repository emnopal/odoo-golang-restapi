package routes

import (
	"log"
	"net/http"

	config "github.com/emnopal/go_postgres/configs"
	controller "github.com/emnopal/go_postgres/controllers"
	"github.com/emnopal/go_postgres/models"
)

func Routes() {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	model := &models.Models{DB: db}
	index := &controller.IndexModels{Model: model}

	http.HandleFunc("/", index.Index)
}
