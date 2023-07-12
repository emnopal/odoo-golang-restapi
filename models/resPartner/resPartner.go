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

func (r *ResPartner) ResPartnerProp() (result *prop.DataProp) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	query := `SELECT reltuples::bigint AS estimate FROM pg_class WHERE oid = 'public.res_partner'::regclass;`

	Length, err := db.Query(query)

	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil
	}

	defer Length.Close()

	var length uint

	for Length.Next() {
		err = Length.Scan(&length)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
		}
	}

	totalPage := uint(math.Floor(float64(length) / float64(r.Limit)))

	result = &prop.DataProp{
		Length:    length,
		TotalPage: totalPage,
	}

	return
}

func (r *ResPartner) GetResPartner(page uint) (result []db_schema.ResPartner, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	query := `SELECT "id", "name", "email", "create_date" FROM "res_partner" ORDER BY "id" offset $1 limit $2`
	offset := page * r.Limit
	Data, err := db.Query(query, offset, r.Limit)

	if err != nil {
		log.Println("Query error occured: ", err.Error())
		return nil, err
	}

	defer Data.Close()

	var RP db_schema.ResPartner

	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		result = append(result, RP)
	}

	return
}
