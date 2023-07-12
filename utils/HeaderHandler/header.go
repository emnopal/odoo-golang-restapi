package utils

import "net/http"

// TODO: Please implement this with more HTTP Header
type HeaderParams struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
	ContentType               string
}

func SetHeader(w http.ResponseWriter, params *HeaderParams) {

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
	w.Header().Set("Access-Control-Allow-Origin", params.AccessControlAllowOrigin)

	// to set which methods is allowed to access this rest api
	w.Header().Set("Access-Control-Allow-Methods", params.AccessControlAllowMethods)

	// to set which headers is allowed to access this rest api
	w.Header().Set("Access-Control-Allow-Headers", params.AccessControlAllowHeaders)

	// to set content type of header
	w.Header().Set("Content-Type", params.ContentType)
}
