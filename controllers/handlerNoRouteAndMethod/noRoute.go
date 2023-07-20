package controllers

import (
	"errors"
	"net/http"

	send "github.com/emnopal/odoo-golang-restapi/models/jsonResponse"
	"github.com/gin-gonic/gin"
)

type NoRouteController struct{}

func (attr *NoRouteController) NoRouteHandler(c *gin.Context) {
	j := &send.JsonSendGetHandler{
		GinContext:            c,
		CustomErrorLogMsg:     "URL is not found",
		CustomErrorRespMsg:    "URL is not found",
		CustomErrorRespStatus: http.StatusNotFound,
	}
	j.SendJsonGet(nil, errors.New(""))
}
