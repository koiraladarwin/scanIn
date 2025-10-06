package models

import "github.com/google/uuid"

type AttendeeCategory struct {
	ID          string  `json:"id"`
	FirebaseID  string  `json:"firebase_id"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	DeletedAt   *string `json:"deleted_at"`
}

type Attendee struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	DeletedAt *string   `json:"deleted_at"`
	AttendeeCategoryID uuid.UUID `json:"attendee_category_id"`
}

type AttendeeRequest struct {
	UserID             uuid.UUID `json:"user_id"`
	TicketID           uuid.UUID `json:"ticket_id"`
	AttendeeCategoryID uuid.UUID `json:"attendee_category_id"`
}
