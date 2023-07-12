package routes

import (
	"net/http"

	controller "github.com/emnopal/go_postgres/controllers"
)

func Routes() {
	index := &controller.IndexController{}
	http.HandleFunc("/", index.Index)
	http.HandleFunc("/contoh", index.Contoh)
}
