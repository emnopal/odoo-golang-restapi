package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	json_schema "github.com/emnopal/go_postgres/schemas/json"
	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (attr *IndexController) Index(c *gin.Context) {

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
