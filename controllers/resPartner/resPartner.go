package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	send "github.com/emnopal/go_postgres/models/jsonResponse"
	resPartner "github.com/emnopal/go_postgres/models/resPartner"
	resPartnerSchema "github.com/emnopal/go_postgres/schemas/db/resPartner"
	"github.com/gin-gonic/gin"
)

type ResPartnerController struct{}

const (
	DefaultLimit             = 100
	DefaultIgnorePerformance = false
)

func (attr *ResPartnerController) GetResPartner(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	pageParams := url.QueryEscape(queryParams.Get("page"))
	limitParams := url.QueryEscape(queryParams.Get("limit"))
	searchParams := url.QueryEscape(queryParams.Get("search"))
	sortParams := url.QueryEscape(queryParams.Get("sort"))
	ignorePerformanceParams := url.QueryEscape(queryParams.Get("ignore_performance"))

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

	RP := &resPartner.ResPartner{}
	if searchParams != "" {
		searchQuery := searchParams
		currentPage := uint(page)
		result, err := RP.GetResPartnerBy(searchQuery, currentPage, uint(limit), sort, ignorePerformance)
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

	currentPage := uint(page)
	result, err := RP.GetResPartner(currentPage, uint(limit), sort, ignorePerformance)

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

	if err := RP.CreateResPartner(&request); err != nil {
		j.SendJsonGet(nil, err)
		return
	}

	j.CustomSuccessLogMsg = fmt.Sprintf("Success send data to: %s", j.GinContext.Request.URL.RequestURI())
	j.CustomSuccessRespMsg = "Success"
	j.CustomSuccessRespStatus = http.StatusCreated

	j.SendJsonGet(nil, nil)
}
