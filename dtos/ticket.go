package dtos

import (
	"github.com/abdullah-mobin/helloBanglaTTS-backend/models"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateTicketDTO struct {
	Subject  string                `json:"subject"`
	Category string                `json:"category"`
	Status   models.TicketStatus   `json:"status,omitempty"`
	Priority models.TicketPriority `json:"priority,omitempty"`
}
type UpdateTicketDTO struct {
	Subject  string                `json:"subject,omitempty"`
	Category string                `json:"category,omitempty"`
	Status   models.TicketStatus   `json:"status,omitempty"`
	Priority models.TicketPriority `json:"priority,omitempty"`
}

func (obj CreateTicketDTO) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Subject, validation.Required),
		validation.Field(&obj.Category, validation.Required),
	)
}
func (obj UpdateTicketDTO) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Status, validation.In(
			models.TicketStatusOpen,
			models.TicketStatusPending,
			models.TicketStatusClose).Error("must be a valid ticket status"),
		),
		validation.Field(&obj.Priority, validation.In(
			models.TicketPriorityHigh,
			models.TicketPriorityMedium,
			models.TicketPriorityLow).Error("must be a valid ticket priority"),
		),
	)
}
