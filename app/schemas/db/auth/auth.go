package schemas

import (
	nulls "github.com/emnopal/odoo-golang-restapi/app/utils/NullHandler"
)

type UserAuth struct {
	ID       nulls.NullString `json:"id"`
	Username nulls.NullString `json:"username"`
	Password nulls.NullString `json:"password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AfterLogin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type User struct {
	ID       nulls.NullString `json:"id"`
	Username nulls.NullString `json:"username"`
	Active   nulls.NullString `json:"is_active"`
	Created  nulls.NullTime   `json:"created"`
	Updated  nulls.NullTime   `json:"updated"`
}

type UserQueryParams struct {
	Page              uint
	Limit             uint
	Search            string
	Sort              string
	IgnorePerformance bool
	MatchExactly      bool
}
