package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	str "github.com/emnopal/go_postgres/utils/Checking"
)

type CommonDB struct{}

func (c *CommonDB) GetTableLength(db *sql.DB, tableName string) (length uint, err error) {
	getLengthDataQuery := `SELECT reltuples::bigint AS estimate FROM pg_class WHERE oid = $1::regclass;`
	tableNameScheme := `public.` + tableName
	Length, err := db.Query(getLengthDataQuery, tableNameScheme)

	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return 0, err
	}

	for Length.Next() {
		err = Length.Scan(&length)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return 0, err
		}
	}
	defer Length.Close()
	return
}

type GetQueryParams struct {
	Columns           string
	Search            map[string]interface{}
	MatchExactly      bool
	Sort              string
	Page              uint
	Limit             uint
	IgnorePerformance bool
}

func (c *CommonDB) StandardGetQuery(db *sql.DB, tableName string, params *GetQueryParams) (Data *sql.Rows, err error) {
	var query string
	page := params.Page
	limit := params.Limit
	columns := params.Columns
	sort := params.Sort
	ignorePerformance := params.IgnorePerformance
	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)
	offset := page * limit // limit the data based on length and limit to improve performance
	if columns == "" {
		columns = "*"
	}

	if ignorePerformance {
		query = fmt.Sprintf(`SELECT %s FROM res_partner ORDER BY %s`, columns, sort)
	} else {
		query = fmt.Sprintf(`SELECT %s FROM res_partner ORDER BY %s OFFSET %d LIMIT %d`, columns, sort, offset, limit)
	}

	Data, err = db.Query(query)
	return
}

func (c *CommonDB) GetQueryById(db *sql.DB, tableName string, params *GetQueryParams) (Data *sql.Rows, err error) {
	columns := params.Columns
	if columns == "" {
		columns = "*"
	}

	query := fmt.Sprintf(`SELECT %s FROM res_partner WHERE %s = %s`, columns, "id", params.Search["id"])

	Data, err = db.Query(query)
	return
}

func (c *CommonDB) GetTableLengthBy(db *sql.DB, tableName string, params *GetQueryParams) (length uint, err error) {
	joinQuery := " UNION "
	if params.MatchExactly {
		joinQuery = " INTERSECT "
	}

	var query string
	var appendQuery []string
	for col, val := range params.Search {
		query = fmt.Sprintf(`(SELECT id FROM %s WHERE %s ILIKE '%%%s%%')`, tableName, col, val)
		appendQuery = append(appendQuery, query)
	}

	finalQuery := strings.Join(appendQuery, joinQuery)
	lengthQuery := fmt.Sprintf(`SELECT COUNT(JOIN_COUNT.id) FROM (%s) JOIN_COUNT`, finalQuery)
	Length, err := db.Query(lengthQuery)

	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return 0, err
	}

	for Length.Next() {
		err = Length.Scan(&length)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return 0, err
		}
	}
	return
}

func (c *CommonDB) GetQueryBy(db *sql.DB, tableName string, params *GetQueryParams) (Data *sql.Rows, err error) {
	joinQuery := " UNION "
	if params.MatchExactly {
		joinQuery = " INTERSECT "
	}

	var query string
	var appendQuery []string
	page := params.Page
	limit := params.Limit
	columns := params.Columns
	sort := params.Sort
	ignorePerformance := params.IgnorePerformance
	offset := page * limit // limit the data based on length and limit to improve performance
	sort = strings.Replace(sort, "+", " ", -1)
	sort = strings.Replace(sort, "%20", " ", -1)
	sort = strings.Replace(sort, "%2C", ",", -1)
	if columns == "" {
		columns = "*"
	}

	for col, val := range params.Search {
		if ignorePerformance {
			query = fmt.Sprintf(`(SELECT %s FROM res_partner WHERE %s ILIKE '%%%s%%' ORDER BY %s)`, columns, col, val, sort)
		} else {
			query = fmt.Sprintf(`(SELECT %s FROM res_partner WHERE %s ILIKE '%%%s%%' ORDER BY %s OFFSET %d LIMIT %d)`, columns, col, val, sort, offset, limit)
		}
		appendQuery = append(appendQuery, query)
	}

	finalQuery := strings.Join(appendQuery, joinQuery)
	Data, err = db.Query(finalQuery)
	return
}

// func (c *CommonDB) StandardPostQuery(db *sql.DB, tableName string, params *GetQueryParams) (Data *sql.Rows, err error) {
// 	return
// }

type PostQueryParams struct {
	Field         map[string]string
	EvaluateField []string
}

func (c *CommonDB) StandardPostQuery(db *sql.DB, tableName string, params *PostQueryParams) (err error) {

	var columnsTmp []string
	var valuesTmp []string

	for col, val := range params.Field {
		if !str.Contains(params.EvaluateField, col) {
			valuesTmp = append(valuesTmp, fmt.Sprintf("'%s'", val))
		} else {
			valuesTmp = append(valuesTmp, val)
		}
		columnsTmp = append(columnsTmp, col)
	}

	columns := strings.Join(columnsTmp, ", ")
	values := strings.Join(valuesTmp, ", ")

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, columns, values)

	tx, err := db.Begin()
	if err != nil {
		log.Println("CreateUsers begin() error: ", err.Error())
		return
	}

	// insert into database
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		log.Println("error insert: ", err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("error commit: ", err.Error())
	}

	return
}
