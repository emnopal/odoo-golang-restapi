package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	json_schema "github.com/emnopal/go_postgres/schemas/json"
	h "github.com/emnopal/go_postgres/utils/HeaderHandler"
)

type IndexController struct{}

func (attr *IndexController) Index(w http.ResponseWriter, req *http.Request) {

	headParams := &h.HeaderParams{
		AccessControlAllowMethods: "GET, POST",
	}
	h.SetHeader(w, headParams)

	message := ""
	status := http.StatusOK

	query_param := req.URL.Query().Get("query_param")

	var t json_schema.ExampleJSON
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
