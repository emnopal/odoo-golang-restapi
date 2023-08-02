package controller

import (
	"fmt"
	"net/http"
	"strconv"

	auth "github.com/emnopal/odoo-golang-restapi/app/models/auth"
	send "github.com/emnopal/odoo-golang-restapi/app/models/jsonResponse"
	authSchema "github.com/emnopal/odoo-golang-restapi/app/schemas/db/auth"
	authConfig "github.com/emnopal/odoo-golang-restapi/app/utils/Token"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (attr *AuthController) Login(c *gin.Context) {
	j := &send.JsonSendGetHandler{GinContext: c}
	authLogin := &auth.Auth{}
	var request authSchema.Login
	if err := c.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			j.CustomErrorLogMsg = "JSON is null"
			j.CustomErrorRespMsg = "JSON is null"
			j.CustomErrorRespStatus = http.StatusBadRequest
		}
		j.SendJsonGet(nil, err)
		return
	}
	// authConfig.DecryptAes256(request, )
	result, err := authLogin.Login(&request)
	if err != nil {
		j.SendJsonGet(nil, err)
		return
	}
	j.CustomSuccessLogMsg = fmt.Sprintf("Success login. Url: %s", j.GinContext.Request.URL.RequestURI())
	j.CustomSuccessRespMsg = j.CustomSuccessLogMsg
	j.CustomSuccessRespStatus = http.StatusOK
	j.SendJsonGet(result, nil)
}

func (attr *AuthController) Profile(c *gin.Context) {
	j := &send.JsonSendGetHandler{GinContext: c}
	id, err := authConfig.ExtractIDFromGin(c)
	if err != nil {
		j.SendJsonGet(nil, err)
		return
	}

	authLogin := &auth.Auth{}
	result, err := authLogin.GetUserById(id)
	if err != nil {
		j.SendJsonGet(nil, err)
		return
	}

	j.SendJsonGet(result, nil)
}

func (attr *AuthController) ProfileBy(c *gin.Context) {
	j := &send.JsonSendGetHandler{GinContext: c}
	param := c.Param("param")

	authLogin := &auth.Auth{}

	_, err := strconv.Atoi(param)
	if err != nil {
		result, err := authLogin.GetUserByUsername(param)
		if err != nil {
			j.SendJsonGet(nil, err)
			return
		}
		j.SendJsonGet(result, nil)
	} else {
		result, err := authLogin.GetUserById(param)
		if err != nil {
			j.SendJsonGet(nil, err)
			return
		}
		j.SendJsonGet(result, nil)
	}

}
