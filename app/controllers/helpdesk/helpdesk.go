package controller

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	helpdesk "github.com/emnopal/odoo-golang-restapi/app/models/helpdesk"
	send "github.com/emnopal/odoo-golang-restapi/app/models/jsonResponse"
	helpdeskSchema "github.com/emnopal/odoo-golang-restapi/app/schemas/db/helpdesk"
	langConfig "github.com/emnopal/odoo-golang-restapi/app/utils/Language"
	"github.com/gin-gonic/gin"
)

type HelpdeskController struct{}

const (
	DefaultLimit             = 100
	DefaultIgnorePerformance = false
	DefaultMatchExcatly      = false
)

func (attr *HelpdeskController) HelpdeskTicketQueryParamsHandler(c *gin.Context) (params *helpdeskSchema.HelpdeskQueryParams) {
	queryParams := c.Request.URL.Query()
	pageParams := url.QueryEscape(queryParams.Get("page"))
	limitParams := url.QueryEscape(queryParams.Get("limit"))
	searchParams := url.QueryEscape(queryParams.Get("search"))
	sortParams := url.QueryEscape(queryParams.Get("sort"))
	ignorePerformanceParams := url.QueryEscape(queryParams.Get("ignore_performance"))
	matchExactlyParams := url.QueryEscape(queryParams.Get("match_exactly"))
	lang := langConfig.GetLangFromHeader(c)

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

	return &helpdeskSchema.HelpdeskQueryParams{
		Page:              uint(page),
		Limit:             uint(limit),
		Search:            searchParams,
		Sort:              sort,
		IgnorePerformance: ignorePerformance,
		MatchExactly:      matchExactly,
		Lang:              lang,
	}
}

func (attr *HelpdeskController) GetHelpdeskTicket(c *gin.Context) {
	params := attr.HelpdeskTicketQueryParamsHandler(c)
	Helpdesk := &helpdesk.Helpdesk{}

	if params.Search != "" {
		result, err := Helpdesk.GetHelpdeskBy(params)
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

	result, err := Helpdesk.GetHelpdeskTicket(params)
	j := &send.JsonSendGetHandler{GinContext: c}
	if err != nil {
		if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
	}
	j.SendJsonGet(result, err)
}

func (attr *HelpdeskController) GetHelpdeskTicketById(c *gin.Context) {
	Helpdesk := &helpdesk.Helpdesk{}
	id := c.Param("id")
	result, err := Helpdesk.GetHelpdeskTicketFromId(id, c)
	j := &send.JsonSendGetHandler{GinContext: c}
	if err != nil {
		if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
	}
	j.SendJsonGet(result, err)
}

func (attr *HelpdeskController) GetHelpdeskTicketStage(c *gin.Context) {
	params := attr.HelpdeskTicketQueryParamsHandler(c)
	params.Sort = "sequence"
	Helpdesk := &helpdesk.Helpdesk{}

	if params.Search != "" {
		result, err := Helpdesk.GetHelpdeskTicketStageBy(params)
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
	result, err := Helpdesk.GetHelpdeskTicketStage(params)
	j := &send.JsonSendGetHandler{GinContext: c}
	if err != nil {
		if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
	}
	j.SendJsonGet(result, err)
}

func (attr *HelpdeskController) GetHelpdeskTicketStageById(c *gin.Context) {
	Helpdesk := &helpdesk.Helpdesk{}
	id := c.Param("id")
	result, err := Helpdesk.GetHelpdeskTicketStageFromId(id, c)
	j := &send.JsonSendGetHandler{GinContext: c}
	if err != nil {
		if err.Error() == "404" || err.Error() == "null" || err.Error() == "null result" {
			err = errors.New("null result")
			j.CustomErrorRespStatus = http.StatusNotFound
		}
	}
	j.SendJsonGet(result, err)
}
