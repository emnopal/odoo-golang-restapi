package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	send "github.com/emnopal/odoo-golang-restapi/app/models/jsonResponse"
	json_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/json"
	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (attr *IndexController) Contoh(c *gin.Context) {

	message := ""
	status := http.StatusOK

	query_param := c.Request.URL.Query().Get("query_param")

	var t json_schema.ExampleJSON
	JSONBody := ""

	if c.Request.Method == "POST" {
		err := json.NewDecoder(c.Request.Body).Decode(&t)
		if err != nil {
			log.Print("WARNING! JSON is empty")
		} else {
			JSONBody = t.JSONBody
		}
	}

	switch c.Request.Method {
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
		message = fmt.Sprintf(`{"message": "Method %s not allowed"}`, c.Request.Method)
		status = http.StatusMethodNotAllowed
	}

	c.JSON(status, message)
}

func (attr *IndexController) Index(c *gin.Context) {

	var t json_schema.ExampleJSON
	var JSONBody string
	err := json.NewDecoder(c.Request.Body).Decode(&t)
	if err != nil {
		log.Print("WARNING! JSON is empty")
		JSONBody = "Empty"
	} else {
		JSONBody = t.JSONBody
	}

	jsonResult := map[string]string{
		"message": "Hello",
		"status":  fmt.Sprintf("%d", http.StatusOK),
		"json":    JSONBody,
		"params":  c.Request.URL.Query().Encode(),
	}

	j := &send.JsonSendGetHandler{GinContext: c}
	j.CustomSuccessRespStatus = http.StatusOK
	j.SendJsonGet(jsonResult, nil)

}
