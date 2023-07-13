package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	send "github.com/emnopal/go_postgres/models/jsonResponse"
	resPartner "github.com/emnopal/go_postgres/models/resPartner"
	"github.com/gin-gonic/gin"
)

type ResPartnerController struct{}

func (attr *ResPartnerController) GetResPartner(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	pageParams := queryParams.Get("page")
	limitParams := queryParams.Get("limit")
	searchParams := queryParams.Get("search")
	ignorePerformanceParams := queryParams.Get("ignore_performance")

	page := 0
	if pageParams != "" {
		var convErr error
		page, convErr = strconv.Atoi(pageParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			page = 0
		}
	}

	limit := 10
	if limitParams != "" {
		var convErr error
		limit, convErr = strconv.Atoi(limitParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			limit = 10
		}
	}

	ignorePerformance := false
	if ignorePerformanceParams != "" {
		var convErr error
		ignorePerformance, convErr = strconv.ParseBool(ignorePerformanceParams)
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			ignorePerformance = false
			ignorePerformanceParams = "false"
		}
	}

	RP := &resPartner.ResPartner{Limit: uint(limit), IgnorePerformance: ignorePerformance}

	if searchParams != "" {
		searchQuery := searchParams
		currentPage := uint(page)
		result, err := RP.GetResPartnerBy(searchQuery, currentPage)
		j := &send.JsonSendGetHandler{GinContext: c}
		if result == nil {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
		j.SendJsonGet(result, err)
		return
	}

	currentPage := uint(page)
	result, err := RP.GetResPartner(currentPage)
	j := &send.JsonSendGetHandler{GinContext: c}
	if result == nil {
		err = errors.New("null result")
		j.CustomErrorRespStatus = http.StatusNotFound
	}
	j.SendJsonGet(result, err)
}

func (attr *ResPartnerController) GetResPartnerById(c *gin.Context) {
	RP := &resPartner.ResPartner{}
	id := c.Param("id")
	result, err := RP.GetResPartnerById(id)
	j := &send.JsonSendGetHandler{GinContext: c}
	if result == nil {
		err = errors.New("null result")
		j.CustomErrorRespStatus = http.StatusNotFound
	}
	j.SendJsonGet(result, err)
}
