package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketStatus string
type TicketPriority string

const (
	TicketStatusOpen    TicketStatus = "open"
	TicketStatusPending TicketStatus = "pending"
	TicketStatusClose   TicketStatus = "close"
)
const (
	TicketPriorityHigh   TicketPriority = "high"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityLow    TicketPriority = "low"
)

type Ticket struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Serial    string             `bson:"serial" json:"serial"`
	Subject   string             `bson:"subject" json:"subject"`
	Category  string             `bson:"category" json:"category"`
	Status    TicketStatus       `bson:"status,omitempty" json:"status,omitempty"`
	Priority  TicketPriority     `bson:"priority,omitempty" json:"priority,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
