package controller

import (
	"fmt"
	"net/http"

	auth "github.com/emnopal/odoo-golang-restapi/models/auth"
	// authConfig "github.com/emnopal/odoo-golang-restapi/utils/Token"
	send "github.com/emnopal/odoo-golang-restapi/models/jsonResponse"
	authSchema "github.com/emnopal/odoo-golang-restapi/schemas/db/auth"
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
