package schemas

import (
	nulls "github.com/emnopal/go_postgres/utils/NullHandler"
)

type ResPartner struct {
	ID         nulls.NullString `json:"id"`
	Name       nulls.NullString `json:"name"`
	Email      nulls.NullString `json:"email"`
	CreateDate nulls.NullTime   `json:"create_date"`
	// CreateDate nulls.NullString `json:"create_date"`
	// Active     string `json:"active"`
	// Language   string `json:"language"`
	// Timezone   string `json:"timezone"`
	// Type       string `json:"type"`
	// Phone      string `json:"phone"`
}
