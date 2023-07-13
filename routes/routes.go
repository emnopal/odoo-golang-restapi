package routes

import (
	"net/http"

	indexController "github.com/emnopal/go_postgres/controllers"
	resPartnerController "github.com/emnopal/go_postgres/controllers/resPartner"
)

func Routes() {
	index := &indexController.IndexController{}
	resPartner := &resPartnerController.ResPartnerController{}

	http.HandleFunc("/contoh", index.Index)
	http.HandleFunc("/", resPartner.GetResPartner)
}
