package models

import "github.com/google/uuid"

type Attendee struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"attendee_id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	DeletedAt *string   `json:"deleted_at"`
}

type AttendeeWithDetails struct {
	ID       string    `json:"id"`
	UserID   uuid.UUID `json:"attendee_id"`
	Position string    `json:"position"`
  Company  string    `json:"company"`
	TicketID  uuid.UUID `json:"ticket_id"`
	DeletedAt *string   `json:"deleted_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ImageURL  string    `json:"image_url"`
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

type AttendeeEnroll struct {
	UserID     uuid.UUID `json:"attendee_id"`
	EventID    uuid.UUID `json:"event_id"`
	ActivityID uuid.UUID `json:"session_id"`
	TicketID   uuid.UUID `json:"ticket_or_invitee_id"`
	FirebaseID string    `json:"firebase_id"`
}

type EnrolledAttendee struct {
	AutoID               int    `json:"auto_id"`
	AttendeeCategoryName string `json:"attendee_category_name"`
	AttendeeName         string `json:"attendee_name"`
	AttendeeImage        string `json:"attendee_image"`
	EventName            string `json:"event_name"`
	SessionName          string `json:"session_name"`
	TicketName           string `json:"entry"`
	Attendee_type        string `json:"type"`
}
