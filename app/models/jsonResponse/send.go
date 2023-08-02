package models

import (
	"fmt"
	"log"
	"net/http"

	r "github.com/emnopal/odoo-golang-restapi/app/schemas/json/response"
	"github.com/gin-gonic/gin"
)

type JsonSendGetHandler struct {
	GinContext              *gin.Context
	CustomErrorLogMsg       string
	CustomSuccessLogMsg     string
	CustomErrorRespMsg      string
	CustomSuccessRespMsg    string
	CustomErrorRespStatus   int
	CustomSuccessRespStatus int
	CustomErrorRespData     interface{}
	CustomSuccessRespData   interface{}
}

func (j *JsonSendGetHandler) SendJsonGet(result interface{}, err error) {

	var resp r.GetResponse

	if err != nil {

		resp.Success = false

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

		j.GinContext.JSON(resp.Status, resp)
		return
	}

	resp.Success = true

	SuccessLogMsg := fmt.Sprintf("Success get: %s", j.GinContext.Request.URL.RequestURI())
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

	j.GinContext.JSON(resp.Status, resp)
}
