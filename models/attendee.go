package models

import "github.com/google/uuid"

type Attendee struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"attendee_id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	DeletedAt *string   `json:"deleted_at"`
}

type AttendeeRequest struct {
	UserID     uuid.UUID `json:"attendee_id"`
	TicketID   uuid.UUID `json:"ticket_or_invitee_id"`
	FirebaseID string    `json:"firebase_id"`
}

type AttendeeActivity struct {
	ID         string    `json:"id"`
	AttendeeID uuid.UUID `json:"attendee_enroll_event_id"`
	ActivityID uuid.UUID `json:"session_id"`
	DeletedAt  *string   `json:"deleted_at"`
}

type AttendeeActivityRequest struct {
	AttendeeID uuid.UUID `json:"attendee_enroll_event_id"`
	ActivityID uuid.UUID `json:"session_id"`
	FirebaseID string    `json:"firebase_id"`
}
