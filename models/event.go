package models

import (
	"github.com/google/uuid"
	"time"
)

type UserWithRole struct {
	ID         uuid.UUID `json:"id"`
	AttendeeId uuid.UUID `json:"attendee_id"`
	FullName   string    `json:"full_name"`
	Company    string    `json:"company"`
	Position   string    `json:"position"`
	AutoId     int       `json:"auto_id"`
	Image_url  string    `json:"image_url"`
	Role       string    `json:"role"`
}

type EventCategoryRequest struct {
	FirebaseID  string     `json:"firebase_id"`
	Tag         string     `json:"tag"`
	Description string     `json:"description"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type EventCategory struct {
	ID          uuid.UUID  `json:"id"`
	Tag         string     `json:"tag"`
	FirebaseID  string     `json:"firebase_id"`
	Description string     `json:"description"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type Event struct {
	ID                  uuid.UUID `json:"id"`
	FirebaseID          string    `json:"firebase_id"`
	EventCategoryID     uuid.UUID `json:"event_category_id"`
	Name                string    `json:"name"`
	EventOrganizer      string    `json:"event_organizer"`
	Description         string    `json:"description"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Location            string    `json:"location"`
	NumberOfParticipant int       `json:"number_of_participant"`
	NumberOfStaff       int       `json:"number_of_staff"`
	StaffCode           *string   `json:"staff_code"`
	AdminCode           *string   `json:"admin_code"`
}

type EventWithDetails struct {
	ID                uuid.UUID `json:"id"`
	EventCategoryID   uuid.UUID `json:"event_category_id"`
	EventCategoryName string `json:"event_category_name"`
	Name              string    `json:"name"`
	EventOrganizer    string    `json:"event_organizer"`
	Description       string    `json:"description"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Location          string    `json:"location"`
	SessionNames      []string  `json:"session_names"`
	InvitationCount    int       `json:"invitation_count"`
	TicketCount       int       `json:"ticket_count"`
	CheckedInCount    int       `json:"checked_in_count"`
}

type EventCreateRequest struct {
	Name                string    `json:"name"`
	EventCategoryID     uuid.UUID `json:"event_category_id"`
	FirebaseID          string    `json:"firebase_id"`
	EventOrganizer      string    `json:"event_organizer"`
	Description         string    `json:"description"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Location            string    `json:"location"`
	NumberOfParticipant int       `json:"number_of_participant"`
}

type EventModifyRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}

type EventInfo struct {
	Event      Event      `json:"event"`
	Activities []Activity `json:"activities"`
}
