package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	config "github.com/emnopal/go_postgres/configs"
	prop "github.com/emnopal/go_postgres/schemas/db/prop"
	db_schema "github.com/emnopal/go_postgres/schemas/db/resPartner"
	nulls "github.com/emnopal/go_postgres/utils/NullHandler"
)

type ResPartner struct{}

const (
	TableName   = "res_partner"
	DefaultCols = `"id", "name", "email", "create_date", "write_date"`
)

func (r *ResPartner) GetResPartner(page uint, limit uint, sort string, ignorePerformance bool) (result *prop.GetDataProp, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get length of data
	var length uint
	getLengthDataQuery := `SELECT reltuples::bigint AS estimate FROM pg_class WHERE oid = $1::regclass;`
	tableNameScheme := `public.` + TableName
	err = db.QueryRow(getLengthDataQuery, tableNameScheme).Scan(&length)
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
	query := fmt.Sprintf(`SELECT %s FROM %s ORDER BY %s`, DefaultCols, TableName, sort)
	if ignorePerformance {
		Data, err = db.Query(query)
	} else {
		query = query + ` OFFSET $1 LIMIT $2`
		Data, err = db.Query(query, offset, limit)
	}

	// raise error when Data is not available
	if Data == nil {
		err = errors.New("404")
		return nil, err
	}

	// get the data to struct
	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner
	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate, &RP.WriteDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
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
	result = &prop.GetDataProp{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}
	return
}

// get by id must be singleton
func (r *ResPartner) GetResPartnerById(id string) (result *prop.GetDataProp, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get data
	// generate query
	var RP db_schema.ResPartner
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE "id" = $1`, DefaultCols, TableName)

	// since get by id must be singleton, so better to us QueryRow
	err = db.QueryRow(query, id).Scan(
		&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate, &RP.WriteDate)

	// raise error if query not available
	if err != nil {
		err = errors.New("404")
		return nil, err
	}

	// returning the result
	result = &prop.GetDataProp{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      RP,
	}
	return
}

func (r *ResPartner) GetResPartnerBy(searchQuery string, page uint, limit uint, sort string, ignorePerformance bool, matchExactly bool) (result *prop.GetDataProp, err error) {
	// initialize database
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// set kind of join query
	joinQuery := " UNION "
	if matchExactly {
		joinQuery = " INTERSECT "
	}
	colsToFind := []string{"name", "email"}

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
	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner
	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate, &RP.WriteDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
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
	result = &prop.GetDataProp{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}
	return
}

func (r *ResPartner) CreateResPartner(request *db_schema.CreateResPartner) (result *prop.CreateDataProp, err error) {
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		return nil, err
	}
	defer db.Close()

	// generate query
	query := fmt.Sprintf(`
		INSERT INTO %s (
			"name", "email", "create_date", "display_name",
			"lang", "tz", "active", "type", "is_company", "write_date"
		)
		VALUES (
			$1, $2, NOW(), $1,
			'en_US', 'Asia/Tokyo', True, 'contact', False, NOW()
		) RETURNING "id", "create_date"`, TableName)

	// transaction
	// start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("begin() error: ", err.Error())
		return nil, err
	}

	// insert into database
	lastInsertId := 0
	var createDate nulls.NullTime
	err = tx.QueryRow(
		query,
		request.Name.String,
		request.Email.String).Scan(&lastInsertId, &createDate)

	// exec, then rollback if error
	if err != nil {
		tx.Rollback()
		log.Println("error insert: ", err.Error())
		return nil, err
	}

	// commit, then rollback if error
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("error commit: ", err.Error())
	}

	result = &prop.CreateDataProp{
		LastInsertId: uint(lastInsertId),
		CreateDate:   createDate,
		Request:      request,
	}

	return result, nil
}

func (r *ResPartner) UpdateResPartner(id string, request *db_schema.UpdateResPartner) (result *prop.UpdateDataProp, err error) {
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get value from id if json doesn't contains other field to make data not null
	dataBeforeChange, err := r.GetResPartnerById(id)
	if err != nil {
		log.Printf("error occured when get the data by id: %s, error: %s", id, err)
		return nil, err
	}

	if request.Name.String == "" {
		request.Name = dataBeforeChange.Result.(db_schema.ResPartner).Name
	}

	if request.Email.String == "" {
		request.Email = dataBeforeChange.Result.(db_schema.ResPartner).Email
	}

	// generate query
	query := fmt.Sprintf(`
		UPDATE %s SET
			"name" = $1,
			"email" = $2,
			"write_date" = NOW()
		WHERE "id" = $3 RETURNING "write_date"`, TableName)

	// transaction
	// start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("begin() error: ", err.Error())
		return nil, err
	}

	// update into database
	var writeDate nulls.NullTime
	err = tx.QueryRow(
		query,
		request.Name.String,
		request.Email.String,
		id).Scan(&writeDate)

	// exec, then rollback if error
	if err != nil {
		tx.Rollback()
		log.Println("error insert: ", err.Error())
		return nil, err
	}

	// commit, then rollback if error
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("error commit: ", err.Error())
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		log.Println("strconv error occured: ", err.Error())
		id_int = 0
	}

	result = &prop.UpdateDataProp{
		ID:        uint(id_int),
		WriteDate: writeDate,
		Request:   request,
	}

	return result, nil
}

func (r *ResPartner) DeleteResPartner(id string) (result *prop.DeleteDataProp, err error) {
	db, err := config.DBConfig()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	defer db.Close()

	// get value from id if json doesn't contains other field to make data not null
	_, err = r.GetResPartnerById(id)
	if err != nil {
		log.Printf("error occured when get the data by id: %s, error: %s", id, err)
		return nil, err
	}

	// generate query
	query := fmt.Sprintf(`DELETE FROM %s WHERE "id" = $1`, TableName)

	// transaction
	// start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("begin() error: ", err.Error())
		return nil, err
	}

	// update into database
	_, err = tx.Exec(query, id)

	// log.Println(results)

	// exec, then rollback if error
	if err != nil {
		tx.Rollback()
		log.Println("error delete: ", err.Error())
		return nil, err
	}

	// commit, then rollback if error
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("error commit: ", err.Error())
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		log.Println("strconv error occured: ", err.Error())
		id_int = 0
	}

	result = &prop.DeleteDataProp{
		ID: uint(id_int),
	}

	return result, nil
}
