package models

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	config "github.com/emnopal/go_postgres/configs"
	prop "github.com/emnopal/go_postgres/schemas/db/prop"
	db_schema "github.com/emnopal/go_postgres/schemas/db/resPartner"
)

type ResPartner struct {
	Limit             uint
	IgnorePerformance bool
}

func (r *ResPartner) GetResPartner(page uint) (result *prop.DataProp, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	// get length of data
	getLengthDataQuery := `SELECT reltuples::bigint AS estimate FROM pg_class WHERE oid = 'public.res_partner'::regclass;`
	Length, err := db.Query(getLengthDataQuery)

	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	var length uint

	for Length.Next() {
		err = Length.Scan(&length)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
	}

	defer Length.Close()

	// get data
	query := `SELECT "id", "name", "email", "create_date" FROM "res_partner" ORDER BY "id" offset $1 limit $2`
	offset := page * r.Limit // limit the data based on length and limit to improve performance

	var Data *sql.Rows

	if r.IgnorePerformance {
		Data, err = db.Query(query, 0, length)
	} else {
		Data, err = db.Query(query, offset, r.Limit)
	}
	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner

	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
	}

	defer Data.Close()

	totalPage := uint(math.Floor(float64(length) / float64(r.Limit)))
	if r.IgnorePerformance {
		totalPage = 0
		page = 0
	}

	result = &prop.DataProp{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}

	return
}

// get by id must be singleton
func (r *ResPartner) GetResPartnerById(id string) (result *prop.DataProp, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	// get data
	query := `SELECT "id", "name", "email", "create_date" FROM "res_partner" WHERE "id" = $1`

	Data, err := db.Query(query, id)
	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner

	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
	}

	defer Data.Close()

	result = &prop.DataProp{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      dataResult,
	}

	return
}

func (r *ResPartner) GetResPartnerBy(searchQuery string, page uint) (result *prop.DataProp, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	// get length of data
	getLengthDataQuery := `
		SELECT SUM(UNION_COUNT.count) FROM ((SELECT count(1) FROM "res_partner" where "name" ilike $1)
		UNION ALL
		(SELECT count(1)  FROM "res_partner" where "email" ilike $1)) UNION_COUNT
	`
	searchQuery = fmt.Sprintf("%%%s%%", searchQuery)
	Length, err := db.Query(getLengthDataQuery, searchQuery)

	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	var length uint

	for Length.Next() {
		err = Length.Scan(&length)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
	}

	defer Length.Close()

	// get data
	query := `
		(SELECT "id", "name", "email", "create_date" FROM "res_partner" where "name" ilike $1 order by id offset $2 limit $3)
		UNION ALL
		(SELECT "id", "name", "email", "create_date" FROM "res_partner" where "email" ilike $1 order by id offset $2 limit $3)
	`
	offset := page * r.Limit // limit the data based on length and limit to improve performance

	searchQuery = fmt.Sprintf("%%%s%%", searchQuery)
	var Data *sql.Rows
	if r.IgnorePerformance {
		Data, err = db.Query(query, searchQuery, 0, length)
	} else {
		Data, err = db.Query(query, searchQuery, offset, r.Limit)
	}
	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner

	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
	}

	defer Data.Close()

	totalPage := uint(math.Floor(float64(length) / float64(r.Limit)))
	if r.IgnorePerformance {
		totalPage = 0
		page = 0
	}

	result = &prop.DataProp{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}

	return
}
