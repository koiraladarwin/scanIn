package models

import "github.com/google/uuid"

type Attendee struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	DeletedAt *string   `json:"deleted_at"`
}

type AttendeeRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	TicketID uuid.UUID `json:"ticket_id"`
}
