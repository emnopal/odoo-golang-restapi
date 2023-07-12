package models

import (
	"log"

	config "github.com/emnopal/go_postgres/configs"
	db_schema "github.com/emnopal/go_postgres/schemas/db/resPartner"
)

func GetResPartner() (result []db_schema.ResPartner, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	query := `SELECT id, name, email, create_date FROM "res_partner"`

	Data, err := db.Query(query)

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
