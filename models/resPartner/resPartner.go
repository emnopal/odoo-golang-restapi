package models

import (
	"log"
	"math"

	config "github.com/emnopal/go_postgres/configs"
	prop "github.com/emnopal/go_postgres/schemas/db/prop"
	db_schema "github.com/emnopal/go_postgres/schemas/db/resPartner"
)

type ResPartner struct {
	Limit uint
}

func (r *ResPartner) GetResPartnerAll(page uint) (result *prop.DataProp, err error) {
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

	Data, err := db.Query(query, offset, r.Limit)
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

	result = &prop.DataProp{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}

	return
}

func (r *ResPartner) GetResPartnerById(id string, page uint) (result *prop.DataProp, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	// get length of data
	getLengthDataQuery := `SELECT COUNT(1) FROM "res_partner" WHERE "id" = $1`
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
	query := `SELECT "id", "name", "email", "create_date" FROM "res_partner" WHERE "id" = $1 ORDER BY "id" offset $2 limit $3`
	offset := page * r.Limit // limit the data based on length and limit to improve performance

	Data, err := db.Query(query, id, offset, r.Limit)
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

	result = &prop.DataProp{
		Length:      length,
		TotalPage:   totalPage,
		CurrentPage: page,
		Result:      dataResult,
	}

	return
}
