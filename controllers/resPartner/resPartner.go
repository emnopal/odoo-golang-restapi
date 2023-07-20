package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	send "github.com/emnopal/odoo-golang-restapi/models/jsonResponse"
	resPartner "github.com/emnopal/odoo-golang-restapi/models/resPartner"
	resPartnerSchema "github.com/emnopal/odoo-golang-restapi/schemas/db/resPartner"
	"github.com/gin-gonic/gin"
)

type ResPartnerController struct{}

const (
	DefaultLimit             = 100
	DefaultIgnorePerformance = false
	DefaultMatchExcatly      = false
)

func (attr *ResPartnerController) ResPartnerQueryParamsHandler(c *gin.Context) (params *resPartnerSchema.ResPartnerQueryParams) {
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

	return &resPartnerSchema.ResPartnerQueryParams{
		Page:              uint(page),
		Limit:             uint(limit),
		Search:            searchParams,
		Sort:              sort,
		IgnorePerformance: ignorePerformance,
		MatchExactly:      matchExactly,
	}
}

func (attr *ResPartnerController) GetResPartner(c *gin.Context) {
	params := attr.ResPartnerQueryParamsHandler(c)
	RP := &resPartner.ResPartner{}
	if params.Search != "" {
		result, err := RP.GetResPartnerBy(params)
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
	result, err := RP.GetResPartner(params)
	j := &send.JsonSendGetHandler{GinContext: c}
	if err != nil {
		if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
	}
	j.SendJsonGet(result, err)
}

func (attr *ResPartnerController) GetResPartnerById(c *gin.Context) {
	RP := &resPartner.ResPartner{}
	id := c.Param("id")
	result, err := RP.GetResPartnerById(id)
	j := &send.JsonSendGetHandler{GinContext: c}
	if err != nil {
		if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
	}
	j.SendJsonGet(result, err)
}

func (attr *ResPartnerController) CreateResPartner(c *gin.Context) {
	j := &send.JsonSendGetHandler{GinContext: c}
	RP := &resPartner.ResPartner{}
	var request resPartnerSchema.CreateResPartner
	if err := c.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			j.CustomErrorLogMsg = "JSON is null"
			j.CustomErrorRespMsg = "JSON is null"
			j.CustomErrorRespStatus = http.StatusBadRequest
		}
		j.SendJsonGet(nil, err)
		return
	}
	result, err := RP.CreateResPartner(&request)
	if err != nil {
		j.SendJsonGet(nil, err)
		return
	}
	j.CustomSuccessLogMsg = fmt.Sprintf("Success send data to: %s", j.GinContext.Request.URL.RequestURI())
	j.CustomSuccessRespMsg = j.CustomSuccessLogMsg
	j.CustomSuccessRespStatus = http.StatusCreated
	j.SendJsonGet(result, nil)
}

func (attr *ResPartnerController) UpdateResPartner(c *gin.Context) {
	id := c.Param("id")
	j := &send.JsonSendGetHandler{GinContext: c}
	RP := &resPartner.ResPartner{}
	var request resPartnerSchema.UpdateResPartner
	if err := c.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			j.CustomErrorLogMsg = "JSON is null"
			j.CustomErrorRespMsg = "JSON is null"
			j.CustomErrorRespStatus = http.StatusBadRequest
		}
		j.SendJsonGet(nil, err)
		return
	}
	result, err := RP.UpdateResPartner(id, &request)
	if err != nil {
		if err.Error() == "404" {
			j.CustomErrorLogMsg = fmt.Sprintf("id: %s not found", id)
			j.CustomErrorRespMsg = fmt.Sprintf("id: %s not found", id)
			j.CustomErrorRespStatus = http.StatusNotFound
		}
		j.SendJsonGet(nil, err)
		return
	}
	j.CustomSuccessLogMsg = fmt.Sprintf("Success update data to: %s with ID: %d", j.GinContext.Request.URL.RequestURI(), result.ID)
	j.CustomSuccessRespMsg = j.CustomSuccessLogMsg
	j.CustomSuccessRespStatus = http.StatusOK
	j.SendJsonGet(result, nil)
}

func (attr *ResPartnerController) DeleteResPartner(c *gin.Context) {
	id := c.Param("id")
	j := &send.JsonSendGetHandler{GinContext: c}
	RP := &resPartner.ResPartner{}
	result, err := RP.DeleteResPartner(id)
	if err != nil {
		if err.Error() == "404" {
			j.CustomErrorLogMsg = fmt.Sprintf("id: %s not found", id)
			j.CustomErrorRespMsg = fmt.Sprintf("id: %s not found", id)
			j.CustomErrorRespStatus = http.StatusNotFound
		}
		j.SendJsonGet(nil, err)
		return
	}
	j.CustomSuccessLogMsg = fmt.Sprintf("Success delete data: %s with ID: %s", j.GinContext.Request.URL.RequestURI(), id)
	j.CustomSuccessRespMsg = j.CustomSuccessLogMsg
	j.CustomSuccessRespStatus = http.StatusOK
	j.SendJsonGet(result, nil)
}
