package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	send "github.com/emnopal/go_postgres/models/jsonResponse"
	resPartner "github.com/emnopal/go_postgres/models/resPartner"
	resPartnerSchema "github.com/emnopal/go_postgres/schemas/db/resPartner"
	"github.com/gin-gonic/gin"
)

type ResPartnerController struct{}

// Get all res partner column
//
// URL: /?page={}&limit={}&search={}&ignore_performance={}
//
// Params:
//
// - page (uint): cursor for pagination in database
//
// - limit (uint): limit for database query (to enhance performance)
//
// - search (string): search for some string inside query
//
// - ignore_performance (boolean): ignoring for some performance tweak.
// If true, some performance tweak will be ignored.
// Note: i don't recommend to use this parameter
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

// Get res partner based on id
//
// URL: /:id
//
// Params:
//
// - id (uint): unique id for res.partner table
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

// Create res partner
//
// URL: /
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
