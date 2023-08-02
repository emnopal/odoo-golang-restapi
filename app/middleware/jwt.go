package middleware

import (
	"errors"
	"net/http"

	send "github.com/emnopal/odoo-golang-restapi/app/models/jsonResponse"
	jwt "github.com/emnopal/odoo-golang-restapi/app/utils/Token"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := jwt.Validate(c)
		if err != nil {
			unauthorizedMsg := "Unauthorized"
			err := errors.New(unauthorizedMsg)
			j := &send.JsonSendGetHandler{GinContext: c}
			j.CustomErrorLogMsg = unauthorizedMsg
			j.CustomErrorRespMsg = unauthorizedMsg
			j.CustomErrorRespStatus = http.StatusUnauthorized
			j.SendJsonGet(nil, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
