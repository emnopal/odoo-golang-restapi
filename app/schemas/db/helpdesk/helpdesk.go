package schemas

import (
	nulls "github.com/emnopal/odoo-golang-restapi/app/utils/NullHandler"
)

type HelpdeskStage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sequence string `json:"sequence"`
}

type HelpdeskTicketListView struct {
	ID             string           `json:"id"`
	CompanyID      nulls.NullString `json:"company_id"`
	CompanyName    nulls.NullString `json:"company_name"`
	CustomerID     nulls.NullString `json:"customer_id"`
	CustomerName   nulls.NullString `json:"customer_name"`
	TicketNumber   nulls.NullString `json:"ticket_number"`
	TicketIssue    nulls.NullString `json:"ticket_issue"`
	ReportedOn     nulls.NullTime   `json:"reported_on"`
	AssignedToID   nulls.NullString `json:"assigned_to_id"`
	AssignedToName nulls.NullString `json:"assigned_to_name"`
	StageID        nulls.NullString `json:"stage_id"`
	StageName      nulls.NullString `json:"stage_name"`
	CreateDate     nulls.NullTime   `json:"create_date"`
}

type HelpdeskTicketFormView struct {
	ID             string           `json:"id"`
	CompanyID      nulls.NullString `json:"company_id"`
	CompanyName    nulls.NullString `json:"company_name"`
	CustomerID     nulls.NullString `json:"customer_id"`
	CustomerName   nulls.NullString `json:"customer_name"`
	TicketNumber   nulls.NullString `json:"ticket_number"`
	TicketIssue    nulls.NullString `json:"ticket_issue"`
	ReportedOn     nulls.NullTime   `json:"reported_on"`
	AssignedToID   nulls.NullString `json:"assigned_to_id"`
	AssignedToName nulls.NullString `json:"assigned_to_name"`
	StageID        nulls.NullString `json:"stage_id"`
	StageName      nulls.NullString `json:"stage_name"`
	Description    nulls.NullString `json:"description"`
	CreateDate     nulls.NullTime   `json:"create_date"`
}

type HelpdeskTicketMessage struct {
	MessageID             string           `json:"message_id"`
	MessageAuthor         nulls.NullString `json:"message_author"`
	MessageBody           nulls.NullString `json:"message_body"`
	MessagePublished      nulls.NullTime   `json:"message_published"`
	AttachmentName        nulls.NullString `json:"attachment_name"`
	AttachmentType        nulls.NullString `json:"attachment_type"`
	AttachmentMimeType    nulls.NullString `json:"attachment_mimetype"`
	AttachmentDescription nulls.NullString `json:"attachment_description"`
	AttachmentURL         nulls.NullString `json:"attachment_url"`
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
