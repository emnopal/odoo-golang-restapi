package utils

import (
	"github.com/gin-gonic/gin"
)

// TODO: Please implement this with more HTTP Header
type HeaderParams struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
	ContentType               string
}

func SetHeader(c *gin.Context, params *HeaderParams) {

	// set default value
	if params.AccessControlAllowOrigin == "" {
		params.AccessControlAllowOrigin = "*"
	}
	if params.AccessControlAllowMethods == "" {
		params.AccessControlAllowMethods = "*"
	}
	if params.AccessControlAllowHeaders == "" {
		params.AccessControlAllowHeaders = "*"
	}
	if params.ContentType == "" {
		params.ContentType = "application/json"
	}

	// to set which origin can access this rest api
	c.Request.Header.Add("Access-Control-Allow-Origin", params.AccessControlAllowOrigin)

	// to set which methods is allowed to access this rest api
	c.Request.Header.Add("Access-Control-Allow-Methods", params.AccessControlAllowMethods)

	// to set which headers is allowed to access this rest api
	c.Request.Header.Add("Access-Control-Allow-Headers", params.AccessControlAllowHeaders)

	// to set content type of header
	c.Request.Header.Add("Content-Type", params.ContentType)
}
