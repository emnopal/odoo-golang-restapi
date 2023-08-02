package middleware

import (
	"net/http"

	utils "github.com/emnopal/odoo-golang-restapi/app/utils/env"
	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_X_API_KEY = "nbEKnc3E1p5a8wK74PkkmJjAsHVxfJIi"
)

func XApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := utils.GetENV("X_API_KEY", DEFAULT_X_API_KEY)
		apiKey := c.Request.Header.Get("x-api-key")

		if apiKey == "" {
			c.String(http.StatusUnauthorized, "Missing API Key")
			c.Abort()
			return
		}
		tokenApiKey := secret

		if apiKey != tokenApiKey {
			c.String(http.StatusUnauthorized, "Invalid API Key")
			c.Abort()
			return
		}

		c.Next()
	}
}
