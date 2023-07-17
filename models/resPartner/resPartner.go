package models

import (
	"errors"
	"log"
	"math"

	config "github.com/emnopal/go_postgres/configs"
	common "github.com/emnopal/go_postgres/models/commonSQL"
	prop "github.com/emnopal/go_postgres/schemas/db/prop"
	db_schema "github.com/emnopal/go_postgres/schemas/db/resPartner"
)

type ResPartner struct{}

const (
	TableName = "res_partner"
)

func (r *ResPartner) GetResPartner(page uint, limit uint, sort string, ignorePerformance bool) (result *prop.DataProp, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	// get length of data
	cdb := common.CommonDB{}
	length, err := cdb.GetTableLength(db, TableName)

	// get data
	Data, err := cdb.StandardGetQuery(db, TableName, &common.GetQueryParams{
		Columns:           "id, name, email, create_date",
		Sort:              sort,
		Page:              page,
		Limit:             limit,
		IgnorePerformance: ignorePerformance,
	})

	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner

	if Data == nil {
		err = errors.New("404")
		return nil, err
	}

	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
	}

	defer Data.Close()

	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	totalPage := uint(math.Floor(float64(length) / float64(limit)))
	if ignorePerformance {
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
	cdb := common.CommonDB{}
	Data, err := cdb.GetQueryById(db, TableName, &common.GetQueryParams{
		Columns: "id, name, email, create_date",
		Search: map[string]interface{}{
			"id": id,
		},
	})

	var RP db_schema.ResPartner
	var dataResult []db_schema.ResPartner

	if Data == nil {
		log.Println("isnil")
		err = errors.New("404")
		return nil, err
	}

	for Data.Next() {
		err = Data.Scan(&RP.ID, &RP.Name, &RP.Email, &RP.CreateDate)
		if err != nil {
			log.Println("Query error occured: ", err.Error())
			return nil, err
		}
		dataResult = append(dataResult, RP)
	}

	defer Data.Close()

	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	result = &prop.DataProp{
		Length:      1,
		TotalPage:   0,
		CurrentPage: 0,
		Result:      dataResult,
	}

	return
}

func (r *ResPartner) GetResPartnerBy(searchQuery string, page uint, limit uint, sort string, ignorePerformance bool) (result *prop.DataProp, err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	// get length of data
	cdb := common.CommonDB{}

	queryParams := &common.GetQueryParams{
		Columns:           "id, name, email, create_date",
		Page:              page,
		Limit:             limit,
		Sort:              sort,
		IgnorePerformance: ignorePerformance,
		MatchExactly:      false,
		Search: map[string]interface{}{
			"name":  searchQuery,
			"email": searchQuery,
		},
	}

	// get length of data
	length, err := cdb.GetTableLengthBy(db, TableName, queryParams)

	// get data
	Data, err := cdb.GetQueryBy(db, TableName, queryParams)

	if Data == nil {
		err = errors.New("404")
		return nil, err
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

	if len(dataResult) == 0 {
		err = errors.New("404")
		return nil, err
	}

	totalPage := uint(math.Floor(float64(length) / float64(limit)))
	if ignorePerformance {
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

func (r *ResPartner) CreateResPartner(request *db_schema.CreateResPartner) (err error) {
	db, err := config.DBConfig()

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	defer db.Close()

	cdb := common.CommonDB{}

	postData := &common.PostQueryParams{
		Field: map[string]string{
			"name":         request.Name.String,
			"email":        request.Email.String,
			"create_date":  "NOW()",
			"display_name": request.Name.String,
			"lang":         "en_US",
			"tz":           "Asia/Tokyo",
			"active":       "True",
			"type":         "contact",
			"is_company":   "False",
			"write_date":   "NOW()",
		},
		EvaluateField: []string{"create_date", "active", "is_company", "write_date"},
	}

	err = cdb.StandardPostQuery(db, TableName, postData)

	return nil
}
