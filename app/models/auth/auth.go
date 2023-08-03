package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"

	config "github.com/emnopal/odoo-golang-restapi/app/configs"
	jwt_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/auth"
	db_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/db/auth"
	res "github.com/emnopal/odoo-golang-restapi/app/schemas/db/result"
	auth "github.com/emnopal/odoo-golang-restapi/app/utils/Token"
)

type Auth struct{}

const (
	TableName       = "res_users"
	DefaultAuthCols = `"id", "login", "password"`
	DefaultCols     = `"id", "login", "active", "create_date", "write_date"`
)

// must be singleton
func (a *Auth) Login(login *db_schema.Login) (after *db_schema.AfterLogin, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get data
	// generate query
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE "login" = $1 AND "active" = True `, DefaultAuthCols, TableName)

	// since get by id must be singleton, so better to us QueryRow
	var U db_schema.UserAuth
	err = db.QueryRow(query, login.Username).Scan(
		&U.ID, &U.Username, &U.Password)

	// raise error if query not available
	if err != nil {
		err = errors.New("404")
		return nil, err
	}

	_, err = auth.Pbkdf2Decoder(login.Password, U.Password.String)

	if err != nil {
		err = errors.New("401")
		return nil, err
	}

	token, err := auth.Generate(&jwt_schema.JWTAccessClaims{
		ID:       U.ID.String,
		Username: U.Username.String,
	})

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		err = errors.New("401")
		return nil, err
	}

	// returning the result
	after = &db_schema.AfterLogin{
		ID:       U.ID.String,
		Username: login.Username,
		Token:    token,
	}
	return
}

func (a *Auth) GetUserBy(params *db_schema.UserQueryParams) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	matchExactly := params.MatchExactly
	searchQuery := params.Search
	sort := params.Sort
	page := params.Page
	limit := params.Limit
	ignorePerformance := params.IgnorePerformance

	// set kind of join query
	joinQuery := " UNION "
	if matchExactly {
		joinQuery = " INTERSECT "
	}
	colsToFind := []string{"login"}

	// length of data
	var appendLengthQueryString string
	var appendLengthQuery []string
	for _, col := range colsToFind {
		appendLengthQueryString = fmt.Sprintf(`(SELECT "id" FROM %s WHERE "%s" ILIKE '%%' || $1 || '%%')`, TableName, col)
		appendLengthQuery = append(appendLengthQuery, appendLengthQueryString)
	}
	lengthQuery := strings.Join(appendLengthQuery, joinQuery)
	finalLengthQuery := fmt.Sprintf(`SELECT COUNT(JOIN_COUNT.id) FROM (%s) JOIN_COUNT`, lengthQuery)

	var length uint
	err = db.QueryRow(finalLengthQuery, searchQuery).Scan(&length)
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
	var query string
	var appendQuery []string
	for _, col := range colsToFind {
		query = fmt.Sprintf(`(SELECT %s FROM %s WHERE %s ILIKE '%%' || $1 || '%%' ORDER BY "%s"`, DefaultCols, TableName, col, sort)
		if !ignorePerformance {
			query = query + ` OFFSET $2 LIMIT $3`
		}
		query = query + `)`
		appendQuery = append(appendQuery, query)
	}

	// querying to db
	finalQuery := strings.Join(appendQuery, joinQuery)
	var Data *sql.Rows
	if ignorePerformance {
		Data, err = db.Query(finalQuery, searchQuery)
	} else {
		Data, err = db.Query(finalQuery, searchQuery, offset, limit)
	}

	// raise error when Data is not available
	if Data == nil {
		err = errors.New("404")
		return nil, err
	}

	// parsing the data
	var U db_schema.User
	var dataResult []db_schema.User
	for Data.Next() {
		err = Data.Scan(&U.ID, &U.Username, &U.Active, &U.Created, &U.Updated)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, U)
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

func (a *Auth) GetUserById(id string) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get data
	// generate query
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE "id" = $1`, DefaultCols, TableName)

	// since get by id must be singleton, so better to us QueryRow
	var U db_schema.User
	err = db.QueryRow(query, id).Scan(
		&U.ID, &U.Username, &U.Active, &U.Created, &U.Updated)

	// raise error if query not available
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		err = errors.New("404")
		return nil, err
	}

	result = &res.ResultGetData{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      U,
	}
	return
}

func (a *Auth) GetUserByUsername(username string) (result *res.ResultGetData, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get data
	// generate query
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE "login" = $1`, DefaultCols, TableName)

	// since get by id must be singleton, so better to us QueryRow
	var U db_schema.User
	err = db.QueryRow(query, username).Scan(
		&U.ID, &U.Username, &U.Active, &U.Created, &U.Updated)

	// raise error if query not available
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		err = errors.New("404")
		return nil, err
	}

	result = &res.ResultGetData{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      U,
	}
	return
}
