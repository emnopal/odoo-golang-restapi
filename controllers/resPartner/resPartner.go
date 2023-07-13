package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	send "github.com/emnopal/go_postgres/models/jsonResponse"
	resPartner "github.com/emnopal/go_postgres/models/resPartner"
	h "github.com/emnopal/go_postgres/utils/HeaderHandler"
)

type ResPartnerController struct{}

func (attr *ResPartnerController) GetResPartner(w http.ResponseWriter, req *http.Request) {
	headParams := &h.HeaderParams{
		AccessControlAllowMethods: "GET",
	}
	h.SetHeader(w, headParams)

	queryParams := req.URL.Query()

	if req.Method != "GET" {
		j := &send.JsonSendGetHandler{W: w, Req: req}
		err := errors.New("method not allowed")
		j.CustomErrorRespStatus = http.StatusMethodNotAllowed
		j.SendJsonGet(nil, err)
		return
	}

	page := 0
	if queryParams.Get("page") != "" {
		var convErr error
		page, convErr = strconv.Atoi(req.URL.Query().Get("page"))
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			page = 0
		}
	}

	limit := 10
	if queryParams.Get("limit") != "" {
		var convErr error
		limit, convErr = strconv.Atoi(req.URL.Query().Get("limit"))
		if convErr != nil {
			log.Println("strconv error occured: ", convErr.Error())
			limit = 10
		}
		if limit > 10_000 {
			limit = 100
		}
	}

	RP := &resPartner.ResPartner{Limit: uint(limit)}
	currentPage := uint(page)
	result, err := RP.GetResPartnerAll(currentPage)
	j := &send.JsonSendGetHandler{W: w, Req: req}
	if result == nil {
		err = errors.New("null result")
		j.CustomErrorRespStatus = http.StatusNotFound
	}
	j.SendJsonGet(result, err)
}
