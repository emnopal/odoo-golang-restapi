package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	auth "github.com/emnopal/odoo-golang-restapi/app/models/auth"
	send "github.com/emnopal/odoo-golang-restapi/app/models/jsonResponse"
	authSchema "github.com/emnopal/odoo-golang-restapi/app/schemas/db/auth"
	authConfig "github.com/emnopal/odoo-golang-restapi/app/utils/Token"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

const (
	DefaultLimit             = 100
	DefaultIgnorePerformance = false
	DefaultMatchExcatly      = false
)

func (attr *AuthController) UserQueryParamsHandler(c *gin.Context) (params *authSchema.UserQueryParams) {
	queryParams := c.Request.URL.Query()
	pageParams := url.QueryEscape(queryParams.Get("page"))
	limitParams := url.QueryEscape(queryParams.Get("limit"))
	searchParams := url.QueryEscape(queryParams.Get("search"))
	sortParams := url.QueryEscape(queryParams.Get("sort"))
	ignorePerformanceParams := url.QueryEscape(queryParams.Get("ignore_performance"))
	matchExactlyParams := url.QueryEscape(queryParams.Get("match_exactly"))

	sort := "id"
	if sortParams != "" {
		sort = sortParams
	}

	page := 0
	if pageParams != "" {
		var convErr error
		page, convErr = strconv.Atoi(pageParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			page = 0
		}
	}

	limit := DefaultLimit
	if limitParams != "" {
		var convErr error
		limit, convErr = strconv.Atoi(limitParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			limit = DefaultLimit
		}
	}

	ignorePerformance := DefaultIgnorePerformance
	if ignorePerformanceParams != "" {
		var convErr error
		ignorePerformance, convErr = strconv.ParseBool(ignorePerformanceParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			ignorePerformance = DefaultIgnorePerformance
		}
	}

	matchExactly := DefaultMatchExcatly
	if matchExactlyParams != "" {
		var convErr error
		matchExactly, convErr = strconv.ParseBool(matchExactlyParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			matchExactly = DefaultMatchExcatly
		}
	}

	return &authSchema.UserQueryParams{
		Page:              uint(page),
		Limit:             uint(limit),
		Search:            searchParams,
		Sort:              sort,
		IgnorePerformance: ignorePerformance,
		MatchExactly:      matchExactly,
	}
}

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
	params := attr.UserQueryParamsHandler(c)

	if params.Search != "" {
		authLogin := &auth.Auth{}
		result, err := authLogin.GetUserBy(params)
		j := &send.JsonSendGetHandler{GinContext: c}
		if err != nil {
			if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
				err = errors.New("null result")
				j.CustomErrorRespStatus = http.StatusNotFound
			}
		}
		j.SendJsonGet(result, err)
		return
	}

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
