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
	"github.com/gin-gonic/gin"
)

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
				) attachment_url,
				iat.access_token attachment_access_token
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
			attachment_name, attachment_type, attachment_mimetype, attachment_description, attachment_url, attachment_access_token
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
			&Message.Attachment.AttachmentURL, &Message.Attachment.AttachmentAccessToken,
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
				) attachment_url,
				iat.access_token attachment_access_token
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
			attachment_name, attachment_type, attachment_mimetype, attachment_description, attachment_url, attachment_access_token
		FROM helpdesk_ticket_message WHERE message_id = $2
	`, rawQuery)

	// since get by id must be singleton, so better to use QueryRow
	err = db.QueryRow(query, ticket_id, message_id).Scan(
		&Message.Message.MessageID, &Message.Message.MessageAuthor,
		&Message.Message.MessageBody, &Message.Message.MessagePublished,
		&Message.Attachment.AttachmentID, &Message.Attachment.AttachmentName,
		&Message.Attachment.AttachmentType, &Message.Attachment.AttachmentMimeType,
		&Message.Attachment.AttachmentDescription,
		&Message.Attachment.AttachmentURL, &Message.Attachment.AttachmentAccessToken,
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
				) attachment_url,
				iat.access_token attachment_access_token
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
		"message_author", "message_body", "message_published",
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
		"attachment_mimetype", "attachment_description", "attachment_url", "attachment_access_token",
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
			&Message.Attachment.AttachmentURL, &Message.Attachment.AttachmentAccessToken,
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
