package schemas

import (
	nulls "github.com/emnopal/odoo-golang-restapi/app/utils/NullHandler"
)

type ResultGetData struct {
	Length      uint        `json:"length"`
	TotalPage   uint        `json:"total_page"`
	CurrentPage uint        `json:"current_page"`
	Result      interface{} `json:"result"`
}

type ResultCreateData struct {
	LastInsertId uint           `json:"last_insert_id"`
	CreateDate   nulls.NullTime `json:"create_date"`
	Request      interface{}    `json:"request"`
}

type ResultUpdateData struct {
	ID        uint           `json:"id"`
	WriteDate nulls.NullTime `json:"write_date"`
	Request   interface{}    `json:"request"`
}

type ResultDeleteData struct {
	ID uint `json:"id"`
}
