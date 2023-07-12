package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	prop "github.com/emnopal/go_postgres/schemas/db/prop"
	r "github.com/emnopal/go_postgres/schemas/json/response"
)

type JsonSendGetHandler struct {
	W                       http.ResponseWriter
	Req                     *http.Request
	CustomErrorLogMsg       string
	CustomSuccessLogMsg     string
	CustomErrorRespMsg      string
	CustomSuccessRespMsg    string
	CustomErrorRespStatus   int
	CustomSuccessRespStatus int
	CustomErrorRespData     interface{}
	CustomSuccessRespData   interface{}
	DataProp                *prop.DataProp
}

func (j *JsonSendGetHandler) SendJsonGet(result interface{}, err error) {

	var resp r.GetResponse

	if err != nil {
		ErrorLogMsg := fmt.Sprintf("Error occured: %s", err.Error())
		if j.CustomErrorLogMsg != "" {
			ErrorLogMsg = j.CustomErrorLogMsg
		}
		log.Print(ErrorLogMsg)

		resp.Status = http.StatusBadRequest
		if j.CustomErrorRespStatus != 0 {
			resp.Status = j.CustomErrorRespStatus
		}

		resp.Message = err.Error()
		if j.CustomErrorRespMsg != "" {
			resp.Message = j.CustomErrorRespMsg
		}

		resp.Data = nil
		if j.CustomErrorRespData != nil {
			resp.Data = j.CustomErrorRespData
		}

		j.W.WriteHeader(resp.Status)
		json.NewEncoder(j.W).Encode(resp)
		return
	}

	SuccessLogMsg := fmt.Sprintf("Success get: %s", j.Req.URL.Path)
	if j.CustomSuccessLogMsg != "" {
		SuccessLogMsg = j.CustomSuccessLogMsg
	}
	log.Print(SuccessLogMsg)

	resp.Status = http.StatusOK
	if j.CustomSuccessRespStatus != 0 {
		resp.Status = j.CustomSuccessRespStatus
	}

	resp.Message = "Success"
	if j.CustomSuccessRespMsg != "" {
		resp.Message = j.CustomSuccessRespMsg
	}

	resp.Data = result
	if j.CustomSuccessRespData != nil {
		resp.Data = j.CustomSuccessRespData
	}

	resp.Length = 0
	resp.TotalPage = 0
	resp.CurrentPage = 0

	if j.DataProp != nil {
		resp.Length = j.DataProp.Length
		resp.TotalPage = j.DataProp.TotalPage
		resp.CurrentPage = j.DataProp.CurrentPage
	}

	j.W.WriteHeader(resp.Status)
	json.NewEncoder(j.W).Encode(resp)
}
