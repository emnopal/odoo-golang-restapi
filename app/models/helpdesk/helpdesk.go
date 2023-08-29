package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"

	config "github.com/emnopal/odoo-golang-restapi/app/configs"
	db_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/db/helpdesk"
	res "github.com/emnopal/odoo-golang-restapi/app/schemas/db/result"
	langConfig "github.com/emnopal/odoo-golang-restapi/app/utils/Language"
	"github.com/gin-gonic/gin"
)

type Helpdesk struct{}

func (h *Helpdesk) GetHelpdeskTicket(params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	page := params.Page
	limit := params.Limit
	sort := params.Sort
	ignorePerformance := params.IgnorePerformance
	lang := params.Lang

	rawQuery := `
		WITH res_users_denorm AS (
			SELECT ru.id id, rp.name name FROM res_users ru
			JOIN res_partner rp ON ru.partner_id = rp.id
		), helpdesk_stage_list_view AS (
			SELECT
				t.id id,
				t.client_id company_id,
				tc.name::json->$1 company_name,
				t.partner_id customer_id,
				tp.name customer_name,
				t.number ticket_number,
				t.name ticket_issue,
				t.create_date reported_on,
				t.user_id assigned_to_id,
				tu.name assigned_to_name,
				t.stage_id stage_id,
				ts.name::json->$1 stage_name,
				t.create_date create_date
			FROM helpdesk_ticket t
			LEFT JOIN helpdesk_ticket_client tc ON t.client_id = tc.id
			LEFT JOIN res_partner tp ON t.partner_id = tp.id
			LEFT JOIN res_users_denorm tu ON t.user_id = tu.id
			LEFT JOIN helpdesk_stage ts ON t.stage_id = ts.id
		)
	`

	// get length of data
	var length uint
	getLengthDataQuery := fmt.Sprintf(`%s SELECT COUNT(1) FROM helpdesk_stage_list_view`, rawQuery)
	err = db.QueryRow(getLengthDataQuery, lang).Scan(&length)
	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	// get data
	// validate sort query
	// prevent sql injection by filtering with regex
	sortValidation, err := regexp.Compile("[a-zA-Z0-9 +%_]*$")
	if !sortValidation.MatchString(sort) {
		sort = "id"
	}

	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)

	// limit the data based on length and limit to improve performance
	offset := page * limit

	// build the query
	var Data *sql.Rows
	query := fmt.Sprintf(`
		%s
		SELECT
			id, ticket_number, ticket_issue, reported_on,
			company_id, company_name,
			customer_id, customer_name,
			assigned_to_id, assigned_to_name,
			stage_id, stage_name,
			create_date
		FROM helpdesk_stage_list_view ORDER BY %s
	`, rawQuery, sort)

	if ignorePerformance {
		Data, err = db.Query(query)
	} else {
		query = query + ` OFFSET $2 LIMIT $3`
		Data, err = db.Query(query, lang, offset, limit)
	}

	// raise error when Data is not available
	if Data == nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	// get the data to struct
	var Ticket db_schema.HelpdeskTicketListView
	var dataResult []db_schema.HelpdeskTicketListView
	for Data.Next() {
		err = Data.Scan(
			&Ticket.ID, &Ticket.TicketNumber, &Ticket.TicketIssue, &Ticket.ReportedOn,
			&Ticket.Company.ID, &Ticket.Company.Name,
			&Ticket.Customer.ID, &Ticket.Customer.Name,
			&Ticket.AssignedUser.ID, &Ticket.AssignedUser.Name,
			&Ticket.Stage.ID, &Ticket.Stage.Name,
			&Ticket.CreateDate,
		)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, Ticket)
	}
	defer Data.Close()

	// raise an error when data is not available
	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	totalPage := uint(math.Floor(float64(length) / float64(limit)))
	if ignorePerformance {
		totalPage = 0
		page = 0
	}
	result = &res.ResultGetData{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}
	return
}

// get by id must be singleton
func (h *Helpdesk) GetHelpdeskTicketFromId(id string, c *gin.Context) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	lang := langConfig.GetLangFromHeader(c)

	// get data
	// generate query
	var Ticket db_schema.HelpdeskTicketFormView
	rawQuery := `
		WITH res_users_denorm AS (
			SELECT ru.id id, rp.name name FROM res_users ru
			JOIN res_partner rp ON ru.partner_id = rp.id
		), helpdesk_stage_form_view AS (
			SELECT
				t.id id,
				t.client_id company_id,
				tc.name::json->$1 company_name,
				t.partner_id customer_id,
				tp.name customer_name,
				t.number ticket_number,
				t.name ticket_issue,
				t.create_date reported_on,
				t.user_id assigned_to_id,
				tu.name assigned_to_name,
				t.stage_id stage_id,
				ts.name::json->$1 stage_name,
				t.description description,
				t.create_date create_date
			FROM helpdesk_ticket t
			LEFT JOIN helpdesk_ticket_client tc ON t.client_id = tc.id
			LEFT JOIN res_partner tp ON t.partner_id = tp.id
			LEFT JOIN res_users_denorm tu ON t.user_id = tu.id
			LEFT JOIN helpdesk_stage ts ON t.stage_id = ts.id
		)
	`
	query := fmt.Sprintf(`
		%s
		SELECT
			id, ticket_number, ticket_issue, reported_on, description,
			company_id, company_name,
			customer_id, customer_name,
			assigned_to_id, assigned_to_name,
			stage_id, stage_name,
			create_date
		FROM helpdesk_stage_form_view WHERE id = $2
	`, rawQuery)

	// since get by id must be singleton, so better to us QueryRow
	err = db.QueryRow(query, lang, id).Scan(
		&Ticket.ID, &Ticket.TicketNumber, &Ticket.TicketIssue,
		&Ticket.ReportedOn, &Ticket.TicketDescription,
		&Ticket.Company.ID, &Ticket.Company.Name,
		&Ticket.Customer.ID, &Ticket.Customer.Name,
		&Ticket.AssignedUser.ID, &Ticket.AssignedUser.Name,
		&Ticket.Stage.ID, &Ticket.Stage.Name,
		&Ticket.CreateDate,
	)

	// raise error if query not available
	if err != nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	result = &res.ResultGetData{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      Ticket,
	}
	return
}

func (h *Helpdesk) GetHelpdeskTicketStage(params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	page := params.Page
	limit := params.Limit
	sort := params.Sort
	ignorePerformance := params.IgnorePerformance
	lang := params.Lang

	rawQuery := `
		WITH helpdesk_stage_list_view AS (
			SELECT
				id, name::json->$1 name, sequence
			FROM helpdesk_stage
			WHERE active=True
		)
	`

	// get length of data
	var length uint
	getLengthDataQuery := fmt.Sprintf(`%s SELECT COUNT(1) FROM helpdesk_stage_list_view`, rawQuery)
	err = db.QueryRow(getLengthDataQuery, lang).Scan(&length)
	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	// get data
	// validate sort query
	// prevent sql injection by filtering with regex
	sortValidation, err := regexp.Compile("[a-zA-Z0-9 +%_]*$")
	if !sortValidation.MatchString(sort) {
		sort = "sequence"
	}

	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)

	// limit the data based on length and limit to improve performance
	offset := page * limit

	// build the query
	var Data *sql.Rows
	query := fmt.Sprintf(`
		%s
		SELECT
			id, name, sequence
		FROM helpdesk_stage_list_view ORDER BY %s
	`, rawQuery, sort)

	if ignorePerformance {
		Data, err = db.Query(query)
	} else {
		query = query + ` OFFSET $2 LIMIT $3`
		Data, err = db.Query(query, lang, offset, limit)
	}

	// raise error when Data is not available
	if Data == nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	// get the data to struct
	var Stage db_schema.HelpdeskStage
	var dataResult []db_schema.HelpdeskStage
	for Data.Next() {
		err = Data.Scan(
			&Stage.ID, &Stage.Name, &Stage.Sequence,
		)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, Stage)
	}
	defer Data.Close()

	// raise an error when data is not available
	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	totalPage := uint(math.Floor(float64(length) / float64(limit)))
	if ignorePerformance {
		totalPage = 0
		page = 0
	}
	result = &res.ResultGetData{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}
	return
}

// get by id must be singleton
func (h *Helpdesk) GetHelpdeskTicketStageFromId(id string, c *gin.Context) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	lang := langConfig.GetLangFromHeader(c)

	// get data
	// generate query
	var Stage db_schema.HelpdeskStage
	rawQuery := `
		WITH helpdesk_stage_form_view AS (
			SELECT
				id, name::json->$1 name, sequence
			FROM helpdesk_stage
			WHERE active = True
		)
	`
	query := fmt.Sprintf(`
		%s
		SELECT
			id, name, sequence
		FROM helpdesk_stage_form_view WHERE id = $2
	`, rawQuery)

	// since get by id must be singleton, so better to us QueryRow
	err = db.QueryRow(query, lang, id).Scan(
		&Stage.ID, &Stage.Name, &Stage.Sequence,
	)

	// raise error if query not available
	if err != nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	result = &res.ResultGetData{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      Stage,
	}
	return
}

func (h *Helpdesk) GetHelpdeskBy(params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	matchExactly := params.MatchExactly
	decodedSearchQuery, err := url.QueryUnescape(params.Search)
	if err != nil {
		log.Println("Query unescape error occured: ", err.Error())
		return nil, err
	}
	sort := params.Sort
	page := params.Page
	limit := params.Limit
	ignorePerformance := params.IgnorePerformance
	lang := params.Lang

	// set kind of join query
	joinQuery := " UNION "
	if matchExactly {
		joinQuery = " INTERSECT "
	}

	rawQuery := `
		WITH res_users_denorm AS (
			SELECT ru.id id, rp.name name FROM res_users ru
			JOIN res_partner rp ON ru.partner_id = rp.id
		), helpdesk_stage_form_view AS (
			SELECT
				t.id id,
				t.client_id company_id,
				tc.name::json->$1::text company_name,
				t.partner_id customer_id,
				tp.name customer_name,
				t.number ticket_number,
				t.name ticket_issue,
				t.create_date reported_on,
				t.user_id assigned_to_id,
				tu.name assigned_to_name,
				t.stage_id stage_id,
				ts.name::json->$1::text stage_name,
				t.description description,
				t.create_date create_date
			FROM helpdesk_ticket t
			LEFT JOIN helpdesk_ticket_client tc ON t.client_id = tc.id
			LEFT JOIN res_partner tp ON t.partner_id = tp.id
			LEFT JOIN res_users_denorm tu ON t.user_id = tu.id
			LEFT JOIN helpdesk_stage ts ON t.stage_id = ts.id
		)
	`

	colsToFind := []string{
		"id", "ticket_number", "ticket_issue", "reported_on",
		"company_id", "company_name",
		"customer_id", "customer_name",
		"assigned_to_id", "assigned_to_name",
		"stage_id", "stage_name",
		"create_date", "description",
	}

	// colsStr := strings.Join(cols, ", ")

	// length of data
	var appendLengthQueryString string
	var appendLengthQuery []string
	for _, col := range colsToFind {
		appendLengthQueryString = fmt.Sprintf(`(SELECT "id" FROM helpdesk_stage_form_view WHERE "%s"::TEXT ILIKE '%%' || $2 || '%%')`, col)
		appendLengthQuery = append(appendLengthQuery, appendLengthQueryString)
	}
	lengthQuery := strings.Join(appendLengthQuery, joinQuery)
	finalLengthQuery := fmt.Sprintf(`%s SELECT COUNT(JOIN_COUNT.id) FROM (%s) JOIN_COUNT`, rawQuery, lengthQuery)

	var length uint
	err = db.QueryRow(finalLengthQuery, lang, decodedSearchQuery).Scan(&length)
	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	// get data
	// validate sort query
	// prevent sql injection by filtering with regex
	sortValidation, err := regexp.Compile("[a-zA-Z0-9 +%_]*$")
	if !sortValidation.MatchString(sort) {
		sort = "id"
	}
	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)

	// limit the data based on length and limit to improve performance
	offset := page * limit

	colsToShow := []string{
		"id", "ticket_number", "ticket_issue", "reported_on",
		"company_id", "company_name::TEXT",
		"customer_id", "customer_name",
		"assigned_to_id", "assigned_to_name",
		"stage_id", "stage_name::TEXT",
		"create_date",
	}

	colsToShowStr := strings.Join(colsToShow, ", ")

	// build the query
	var query string
	var appendQuery []string
	for _, col := range colsToFind {
		query = fmt.Sprintf(`(SELECT %s FROM helpdesk_stage_form_view WHERE %s::TEXT ILIKE '%%' || $2 || '%%' ORDER BY "%s"`, colsToShowStr, col, sort)
		if !ignorePerformance {
			query = query + ` OFFSET $3 LIMIT $4`
		}
		query = query + `)`
		appendQuery = append(appendQuery, query)
	}

	// querying to db
	finalQuery := strings.Join(appendQuery, joinQuery)
	finalQuery = rawQuery + finalQuery

	var Data *sql.Rows
	if ignorePerformance {
		Data, err = db.Query(finalQuery, lang, decodedSearchQuery)
	} else {
		Data, err = db.Query(finalQuery, lang, decodedSearchQuery, offset, limit)
	}

	// check if there are errors
	if err != nil {
		log.Println("Query error occured: ", err.Error())
	}

	// raise error when Data is not available
	if Data == nil {
		err = errors.New("404")
		return nil, err
	}

	// get the data to struct
	var Ticket db_schema.HelpdeskTicketListView
	var dataResult []db_schema.HelpdeskTicketListView
	for Data.Next() {
		err = Data.Scan(
			&Ticket.ID, &Ticket.TicketNumber, &Ticket.TicketIssue, &Ticket.ReportedOn,
			&Ticket.Company.ID, &Ticket.Company.Name,
			&Ticket.Customer.ID, &Ticket.Customer.Name,
			&Ticket.AssignedUser.ID, &Ticket.AssignedUser.Name,
			&Ticket.Stage.ID, &Ticket.Stage.Name,
			&Ticket.CreateDate,
		)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, Ticket)
	}
	defer Data.Close()

	// raise an error when data is not available
	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	totalPage := uint(math.Floor(float64(length) / float64(limit)))
	if ignorePerformance {
		totalPage = 0
		page = 0
	}
	result = &res.ResultGetData{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}
	return
}

func (h *Helpdesk) GetHelpdeskTicketStageBy(params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	matchExactly := params.MatchExactly
	decodedSearchQuery, err := url.QueryUnescape(params.Search)
	if err != nil {
		log.Println("Query unescape error occured: ", err.Error())
		return nil, err
	}
	sort := params.Sort
	page := params.Page
	limit := params.Limit
	ignorePerformance := params.IgnorePerformance
	lang := params.Lang

	// set kind of join query
	joinQuery := " UNION "
	if matchExactly {
		joinQuery = " INTERSECT "
	}

	rawQuery := `
			WITH helpdesk_stage_list_view AS (
				SELECT
					id, name::json->$1 name, sequence
				FROM helpdesk_stage
				WHERE active=True
			)
		`

	colsToFind := []string{
		"id", "name", "sequence",
	}

	// colsStr := strings.Join(cols, ", ")

	// length of data
	var appendLengthQueryString string
	var appendLengthQuery []string
	for _, col := range colsToFind {
		appendLengthQueryString = fmt.Sprintf(`(SELECT "id" FROM helpdesk_stage_list_view WHERE "%s"::TEXT ILIKE '%%' || $2 || '%%')`, col)
		appendLengthQuery = append(appendLengthQuery, appendLengthQueryString)
	}
	lengthQuery := strings.Join(appendLengthQuery, joinQuery)
	finalLengthQuery := fmt.Sprintf(`%s SELECT COUNT(JOIN_COUNT.id) FROM (%s) JOIN_COUNT`, rawQuery, lengthQuery)

	var length uint
	err = db.QueryRow(finalLengthQuery, lang, decodedSearchQuery).Scan(&length)
	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	// get data
	// validate sort query
	// prevent sql injection by filtering with regex
	sortValidation, err := regexp.Compile("[a-zA-Z0-9 +%_]*$")
	if !sortValidation.MatchString(sort) {
		sort = "sequence"
	}

	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)

	// limit the data based on length and limit to improve performance
	offset := page * limit

	colsToShow := []string{
		"id", "name::TEXT", "sequence",
	}

	colsToShowStr := strings.Join(colsToShow, ", ")

	// build the query
	var query string
	var appendQuery []string
	for _, col := range colsToFind {
		query = fmt.Sprintf(`(SELECT %s FROM helpdesk_stage_list_view WHERE %s::TEXT ILIKE '%%' || $2 || '%%' ORDER BY "%s"`, colsToShowStr, col, sort)
		if !ignorePerformance {
			query = query + ` OFFSET $3 LIMIT $4`
		}
		query = query + `)`
		appendQuery = append(appendQuery, query)
	}

	// querying to db
	finalQuery := strings.Join(appendQuery, joinQuery)
	finalQuery = rawQuery + finalQuery

	var Data *sql.Rows
	if ignorePerformance {
		Data, err = db.Query(finalQuery, lang, decodedSearchQuery)
	} else {
		Data, err = db.Query(finalQuery, lang, decodedSearchQuery, offset, limit)
	}

	// check if there are errors
	if err != nil {
		log.Println("Query error occured: ", err.Error())
	}

	// raise error when Data is not available
	if Data == nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	// get the data to struct
	var Stage db_schema.HelpdeskStage
	var dataResult []db_schema.HelpdeskStage
	for Data.Next() {
		err = Data.Scan(
			&Stage.ID, &Stage.Name, &Stage.Sequence,
		)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, Stage)
	}
	defer Data.Close()

	// raise an error when data is not available
	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	totalPage := uint(math.Floor(float64(length) / float64(limit)))
	if ignorePerformance {
		totalPage = 0
		page = 0
	}
	result = &res.ResultGetData{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}
	return
}
