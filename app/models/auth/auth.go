package model

import (
	"errors"
	"fmt"
	"log"

	config "github.com/emnopal/odoo-golang-restapi/app/configs"
	db_schema "github.com/emnopal/odoo-golang-restapi/app/schemas/db/auth"
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
	var U db_schema.User
	err = db.QueryRow(query, login.Username).Scan(
		&U.ID, &U.Username, &U.Password)

	// raise error if query not available
	if err != nil {
		err = errors.New("404")
		return nil, err
	}

	_, err = auth.Pbkdf2Decoder(login.Password, U.Password)

	if err != nil {
		err = errors.New("401")
		return nil, err
	}

	token, err := auth.Generate(U.ID)

	if err != nil {
		err = errors.New("401")
		return nil, err
	}

	// returning the result
	after = &db_schema.AfterLogin{
		ID:       U.ID,
		Username: login.Username,
		Token:    token,
	}
	return
}
