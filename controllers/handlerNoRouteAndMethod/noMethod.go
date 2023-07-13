package controllers

import (
	"errors"
	"net/http"

	send "github.com/emnopal/go_postgres/models/jsonResponse"
	"github.com/gin-gonic/gin"
)

type NoMethodController struct{}

func (attr *NoMethodController) NoMethodHandler(c *gin.Context) {
	j := &send.JsonSendGetHandler{
		GinContext:            c,
		CustomErrorLogMsg:     "Method is not allowed",
		CustomErrorRespMsg:    "Method is not allowed",
		CustomErrorRespStatus: http.StatusMethodNotAllowed,
	}
	j.SendJsonGet(nil, errors.New(""))
}
