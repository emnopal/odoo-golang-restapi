package schemas

import (
	nulls "github.com/emnopal/odoo-golang-restapi/app/utils/NullHandler"
)

type HelpdeskStage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sequence string `json:"sequence"`
}

type Company struct {
	ID   nulls.NullString `json:"id"`
	Name nulls.NullString `json:"name"`
}

type Customer struct {
	ID   nulls.NullString `json:"id"`
	Name nulls.NullString `json:"name"`
}

type AssignedUser struct {
	ID   nulls.NullString `json:"id"`
	Name nulls.NullString `json:"name"`
}

type Stage struct {
	ID   nulls.NullString `json:"id"`
	Name nulls.NullString `json:"name"`
}

type Message struct {
	MessageID        string           `json:"message_id"`
	MessageAuthor    nulls.NullString `json:"message_author"`
	MessageBody      nulls.NullString `json:"message_body"`
	MessagePublished nulls.NullTime   `json:"message_published"`
}

type Attachment struct {
	AttachmentID          nulls.NullString `json:"attachment_id"`
	AttachmentName        nulls.NullString `json:"attachment_name"`
	AttachmentType        nulls.NullString `json:"attachment_type"`
	AttachmentMimeType    nulls.NullString `json:"attachment_mimetype"`
	AttachmentDescription nulls.NullString `json:"attachment_description"`
	AttachmentURL         nulls.NullString `json:"attachment_url"`
	AttachmentAccessToken nulls.NullString `json:"attachment_access_token"`
}

type HelpdeskTicketListView struct {
	ID           string           `json:"id"`
	TicketNumber nulls.NullString `json:"ticket_number"`
	TicketIssue  nulls.NullString `json:"ticket_issue"`
	ReportedOn   nulls.NullTime   `json:"reported_on"`
	Company      Company          `json:"company"`
	Customer     Customer         `json:"customer"`
	AssignedUser AssignedUser     `json:"assigned_user"`
	Stage        Stage            `json:"stage"`
}

type HelpdeskTicketFormView struct {
	ID                string           `json:"id"`
	TicketNumber      nulls.NullString `json:"ticket_number"`
	TicketIssue       nulls.NullString `json:"ticket_issue"`
	TicketDescription nulls.NullString `json:"ticket_description"`
	ReportedOn        nulls.NullTime   `json:"reported_on"`
	Company           Company          `json:"company"`
	Customer          Customer         `json:"customer"`
	AssignedUser      AssignedUser     `json:"assigned_user"`
	Stage             Stage            `json:"stage"`
}

type HelpdeskTicketMessage struct {
	Message    Message    `json:"message"`
	Attachment Attachment `json:"attachment"`
}

type HelpdeskQueryParams struct {
	Page              uint
	Limit             uint
	Search            string
	Sort              string
	IgnorePerformance bool
	MatchExactly      bool
	Lang              string
}
