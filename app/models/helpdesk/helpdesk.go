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

// Helpdesk Ticket
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
				ts.name::json->$1 stage_name
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
			stage_id, stage_name
		FROM helpdesk_stage_list_view ORDER BY %s
	`, rawQuery, sort)

	if ignorePerformance {
		Data, err = db.Query(query, lang)
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
				t.description description
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
			stage_id, stage_name
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

func (h *Helpdesk) GetHelpdeskTicketBy(params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
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
				t.description description
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
		"stage_id", "stage_name", "description",
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
	}

	colsToShowStr := strings.Join(colsToShow, ", ")

	// build the query
	var query string
	var appendQuery []string
	for _, col := range colsToFind {
		query = fmt.Sprintf(`(SELECT %s FROM helpdesk_stage_form_view WHERE %s::TEXT ILIKE '%%' || $2 || '%%' ORDER BY %s`, colsToShowStr, col, sort)
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

// Helpdesk Ticket Stage
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
		Data, err = db.Query(query, lang)
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
		query = fmt.Sprintf(`(SELECT %s FROM helpdesk_stage_list_view WHERE %s::TEXT ILIKE '%%' || $2 || '%%' ORDER BY %s`, colsToShowStr, col, sort)
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

// Helpdesk Ticket Message
func (h *Helpdesk) GetHelpdeskTicketMessage(ticket_id string, params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
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
	// lang := params.Lang

	rawQuery := `
		WITH helpdesk_ticket_message AS (
			SELECT
				me.id message_id,
				p.name message_author,
				me.body message_body,
				me.create_date message_published,
				iat.id attachment_id,
				iat.name attachment_name,
				iat.type attachment_type,
				iat.mimetype attachment_mimetype,
				iat.description attachment_description,
				CONCAT(
					'/web/content/',
					iat.id,
					'?download=true&access_token=',
					iat.access_token
				) attachment_url
			FROM
				mail_message me
				LEFT JOIN message_attachment_rel marel ON me.id=marel.message_id
				LEFT JOIN ir_attachment iat ON marel.attachment_id=iat.id
				LEFT JOIN res_partner p ON me.author_id=p.id
			WHERE
				me.model = 'helpdesk.ticket'
				AND me.res_id = $1
		)
	`

	// get length of data
	var length uint
	getLengthDataQuery := fmt.Sprintf(`%s SELECT COUNT(1) FROM helpdesk_ticket_message`, rawQuery)
	err = db.QueryRow(getLengthDataQuery, ticket_id).Scan(&length)
	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	// get data
	// validate sort query
	// prevent sql injection by filtering with regex
	sortValidation, err := regexp.Compile("[a-zA-Z0-9 +%_]*$")
	if !sortValidation.MatchString(sort) {
		sort = "message_published desc"
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
			message_id, message_author, message_body, message_published, attachment_id,
			attachment_name, attachment_type, attachment_mimetype, attachment_description, attachment_url
		FROM helpdesk_ticket_message ORDER BY %s
	`, rawQuery, sort)

	if ignorePerformance {
		Data, err = db.Query(query, ticket_id)
	} else {
		query = query + ` OFFSET $2 LIMIT $3`
		Data, err = db.Query(query, ticket_id, offset, limit)
	}

	// raise error when Data is not available
	if Data == nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	// get the data to struct
	var Message db_schema.HelpdeskTicketMessage
	var dataResult []db_schema.HelpdeskTicketMessage
	for Data.Next() {
		err = Data.Scan(
			&Message.Message.MessageID, &Message.Message.MessageAuthor,
			&Message.Message.MessageBody, &Message.Message.MessagePublished,
			&Message.Attachment.AttachmentID, &Message.Attachment.AttachmentName,
			&Message.Attachment.AttachmentType, &Message.Attachment.AttachmentMimeType,
			&Message.Attachment.AttachmentDescription,
			&Message.Attachment.AttachmentURL,
		)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, Message)
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
func (h *Helpdesk) GetHelpdeskTicketMessageFromId(ticket_id string, message_id string, c *gin.Context) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	var Message db_schema.HelpdeskTicketMessage
	rawQuery := `
		WITH helpdesk_ticket_message AS (
			SELECT
				me.id message_id,
				p.name message_author,
				me.body message_body,
				me.create_date message_published,
				iat.id attachment_id,
				iat.name attachment_name,
				iat.type attachment_type,
				iat.mimetype attachment_mimetype,
				iat.description attachment_description,
				CONCAT(
					'/web/content/',
					iat.id,
					'?download=true&access_token=',
					iat.access_token
				) attachment_url
			FROM
				mail_message me
				LEFT JOIN message_attachment_rel marel ON me.id=marel.message_id
				LEFT JOIN ir_attachment iat ON marel.attachment_id=iat.id
				LEFT JOIN res_partner p ON me.author_id=p.id
			WHERE
				me.model = 'helpdesk.ticket'
				AND me.res_id = $1
		)
	`
	query := fmt.Sprintf(`
		%s
		SELECT
			message_id, message_author, message_body, message_published, attachment_id,
			attachment_name, attachment_type, attachment_mimetype, attachment_description, attachment_url
		FROM helpdesk_ticket_message WHERE message_id = $2
	`, rawQuery)

	// since get by id must be singleton, so better to use QueryRow
	err = db.QueryRow(query, ticket_id, message_id).Scan(
		&Message.Message.MessageID, &Message.Message.MessageAuthor,
		&Message.Message.MessageBody, &Message.Message.MessagePublished,
		&Message.Attachment.AttachmentID, &Message.Attachment.AttachmentName,
		&Message.Attachment.AttachmentType, &Message.Attachment.AttachmentMimeType,
		&Message.Attachment.AttachmentDescription,
		&Message.Attachment.AttachmentURL,
	)

	// raise error if query not available
	if err != nil {
		log.Println(err)
		err = errors.New("404")
		return nil, err
	}

	result = &res.ResultGetData{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      Message,
	}
	return
}

func (h *Helpdesk) GetHelpdeskTicketMessageBy(ticket_id string, params *db_schema.HelpdeskQueryParams) (result *res.ResultGetData, err error) {
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
	// lang := params.Lang

	// set kind of join query
	joinQuery := " UNION "
	if matchExactly {
		joinQuery = " INTERSECT "
	}

	rawQuery := `
		WITH helpdesk_ticket_message AS (
			SELECT
				me.id message_id,
				p.name message_author,
				me.body message_body,
				me.create_date message_published,
				iat.id attachment_id,
				iat.name attachment_name,
				iat.type attachment_type,
				iat.mimetype attachment_mimetype,
				iat.description attachment_description,
				CONCAT(
					'/web/content/',
					iat.id,
					'?download=true&access_token=',
					iat.access_token
				) attachment_url
			FROM
				mail_message me
				LEFT JOIN message_attachment_rel marel ON me.id=marel.message_id
				LEFT JOIN ir_attachment iat ON marel.attachment_id=iat.id
				LEFT JOIN res_partner p ON me.author_id=p.id
			WHERE
				me.model = 'helpdesk.ticket'
				AND me.res_id = $1
		)
	`

	colsToFind := []string{
		"message_id", "message_author", "message_body",
		"message_published", "attachment_id", "attachment_name", "attachment_type",
		"attachment_mimetype", "attachment_description", "attachment_url",
	}

	// colsStr := strings.Join(cols, ", ")

	// length of data
	var appendLengthQueryString string
	var appendLengthQuery []string
	for _, col := range colsToFind {
		appendLengthQueryString = fmt.Sprintf(`(SELECT "message_id" FROM helpdesk_ticket_message WHERE "%s"::TEXT ILIKE '%%' || $2 || '%%')`, col)
		appendLengthQuery = append(appendLengthQuery, appendLengthQueryString)
	}
	lengthQuery := strings.Join(appendLengthQuery, joinQuery)
	finalLengthQuery := fmt.Sprintf(`%s SELECT COUNT(JOIN_COUNT.message_id) FROM (%s) JOIN_COUNT`, rawQuery, lengthQuery)

	var length uint
	err = db.QueryRow(finalLengthQuery, ticket_id, decodedSearchQuery).Scan(&length)
	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	// get data
	// validate sort query
	// prevent sql injection by filtering with regex
	sortValidation, err := regexp.Compile("[a-zA-Z0-9 +%_]*$")
	if !sortValidation.MatchString(sort) {
		sort = "message_published desc"
	}

	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)

	// limit the data based on length and limit to improve performance
	offset := page * limit

	colsToShow := []string{
		"message_id", "message_author", "message_body",
		"message_published", "attachment_id", "attachment_name", "attachment_type",
		"attachment_mimetype", "attachment_description", "attachment_url",
	}

	colsToShowStr := strings.Join(colsToShow, ", ")

	// build the query
	var query string
	var appendQuery []string
	for _, col := range colsToFind {
		query = fmt.Sprintf(`(SELECT %s FROM helpdesk_ticket_message WHERE %s::TEXT ILIKE '%%' || $2 || '%%' ORDER BY %s`, colsToShowStr, col, sort)
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
		Data, err = db.Query(finalQuery, ticket_id, decodedSearchQuery)
	} else {
		Data, err = db.Query(finalQuery, ticket_id, decodedSearchQuery, offset, limit)
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
	var Message db_schema.HelpdeskTicketMessage
	var dataResult []db_schema.HelpdeskTicketMessage
	for Data.Next() {
		err = Data.Scan(
			&Message.Message.MessageID, &Message.Message.MessageAuthor,
			&Message.Message.MessageBody, &Message.Message.MessagePublished,
			&Message.Attachment.AttachmentID, &Message.Attachment.AttachmentName,
			&Message.Attachment.AttachmentType, &Message.Attachment.AttachmentMimeType,
			&Message.Attachment.AttachmentDescription,
			&Message.Attachment.AttachmentURL,
		)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, Message)
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
