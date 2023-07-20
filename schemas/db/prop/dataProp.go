package schemas

import (
	nulls "github.com/emnopal/go_postgres/utils/NullHandler"
)

type GetDataProp struct {
	Length      uint        `json:"length"`
	TotalPage   uint        `json:"total_page"`
	CurrentPage uint        `json:"current_page"`
	Result      interface{} `json:"result"`
}

type CreateDataProp struct {
	LastInsertId uint           `json:"last_insert_id"`
	CreateDate   nulls.NullTime `json:"create_date"`
	Request      interface{}    `json:"request"`
}

type UpdateDataProp struct {
	ID        uint           `json:"id"`
	WriteDate nulls.NullTime `json:"write_date"`
	Request   interface{}    `json:"request"`
}

type DeleteDataProp struct {
	ID uint `json:"id"`
}
