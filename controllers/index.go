package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	send "github.com/emnopal/go_postgres/models/jsonResponse"
	resPartner "github.com/emnopal/go_postgres/models/resPartner"
	json_schema "github.com/emnopal/go_postgres/schemas/json"
	h "github.com/emnopal/go_postgres/utils/HeaderHandler"
)

type IndexController struct{}

func (attr *IndexController) Index(w http.ResponseWriter, req *http.Request) {
	headParams := &h.HeaderParams{}
	h.SetHeader(w, headParams)

	queryParams := req.URL.Query()

	if req.Method != "GET" {
		j := &send.JsonSendGetHandler{W: w, Req: req}
		err := errors.New("method not allowed")
		j.CustomErrorRespStatus = http.StatusMethodNotAllowed
		j.SendJsonGet(nil, err)
		return
	}

	page := 0
	if queryParams.Get("page") != "" {
		var convErr error
		page, convErr = strconv.Atoi(req.URL.Query().Get("page"))
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			page = 0
		}
	}

	RP := &resPartner.ResPartner{Limit: 30}
	resultLength := RP.ResPartnerProp()
	resultLength.CurrentPage = uint(page)
	result, err := RP.GetResPartner(resultLength.CurrentPage)
	j := &send.JsonSendGetHandler{W: w, Req: req, DataProp: resultLength}
	if result == nil {
		err = errors.New("null result")
		j.CustomErrorRespStatus = http.StatusNotFound
	}
	j.SendJsonGet(result, err)
}

func (attr *IndexController) Contoh(w http.ResponseWriter, req *http.Request) {

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
