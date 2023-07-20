package schemas

import (
	nulls "github.com/emnopal/go_postgres/utils/NullHandler"
)

type ResPartner struct {
	ID         string           `json:"id"`
	Name       nulls.NullString `json:"name"`
	Email      nulls.NullString `json:"email"`
	CreateDate nulls.NullTime   `json:"create_date"`
	WriteDate  nulls.NullTime   `json:"write_date"`
	// Active     string `json:"active"`
	// Language   string `json:"language"`
	// Timezone   string `json:"timezone"`
	// Type       string `json:"type"`
	// Phone      string `json:"phone"`
}

type CreateResPartner struct {
	Name  nulls.NullString `json:"name"`
	Email nulls.NullString `json:"email"`
}

type UpdateResPartner struct {
	Name  nulls.NullString `json:"name"`
	Email nulls.NullString `json:"email"`
}
