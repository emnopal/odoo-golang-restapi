package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	helper "github.com/emnopal/go_helper"
	"github.com/emnopal/go_postgres/models"
	schemas "github.com/emnopal/go_postgres/schemas/json"
)

type IndexModels struct {
	Model *models.Models
}

func (m *IndexModels) Index(w http.ResponseWriter, req *http.Request) {

	log.Print(m.Model.DB)

	headParams := &helper.HeaderParams{
		AccessControlAllowMethods: "GET, POST",
	}
	helper.SetHeader(w, headParams)

	message := ""
	status := http.StatusOK

	query_param := req.URL.Query().Get("query_param")

	var t schemas.ExampleJSON
	JSONBody := ""

	if req.Method == "POST" {
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			log.Print("WARNING! JSON is empty")
		} else {
			JSONBody = t.JSONBody
		}
	}

	switch req.Method {
	case "GET":
		message = fmt.Sprintf(`{
			"message": "Hello World from GET",
			"query_params": "%s"
		}`, query_param)
	case "POST":
		message = fmt.Sprintf(`{
			"message": "Hello World from POST",
			"json_body": "%s"
		}`, JSONBody)
	default:
		message = fmt.Sprintf(`{"message": "Method %s not allowed"}`, req.Method)
		status = http.StatusMethodNotAllowed
	}

	w.WriteHeader(status)
	w.Write([]byte(message))
}
