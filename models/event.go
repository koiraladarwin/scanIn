package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered person in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Company   string    `json:"company"`
	Position  string    `json:"position"`
	Image_url string    `json:"image_url"`
	AutoId    int       `json:"auto_id"`
	Role      string    `json:"role"` // "participant", "staff", "member"
}

// User represents a registered person in the system
type UserRequest struct {
	FullName  string `json:"full_name"`
	Company   string `json:"company"`
	Position  string `json:"position"`
	Image_url string `json:"image_url"`
	Role      string `json:"role"` // "participant", "staff", "member"
}

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

// Event represents a specific event like "Teej 2025"
type Event struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Location            string    `json:"location"`
	NumberOfParticipant int       `json:"number_of_participant"`
}

// Activity is a feature or sub-event within an event (e.g. Registration, Food)
type Activity struct {
	ID                  uuid.UUID `json:"id"`
	EventID             uuid.UUID `json:"event_id"` // FK to Event
	Name                string    `json:"name"`     // e.g. Registration, Dinner
	Type                string    `json:"type"`     // e.g. "registration", "food", etc.
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	NumberOfScanedUsers int       `json:"number_of_scaned_users"`
}

// Attendee links a user to an event
type Attendee struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`  // FK to User
	EventID uuid.UUID `json:"event_id"` // FK to Event
}

// CheckInLog is created when an attendee scans at an activity
type CheckInLog struct {
	ID         uuid.UUID `json:"id"`
	AttendeeID uuid.UUID `json:"attendee_id"` // FK to Attendee
	ActivityID uuid.UUID `json:"activity_id"` // FK to Activity
	ScannedAt  time.Time `json:"scanned_at"`
	Status     string    `json:"status"`     // e.g. "success", "duplicate", "invalid"
	ScannedBy  string    `json:"scanned_by"` // FK to User (staff)
}

type CheckInRespose struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"full_name"`
	AttendeeID   uuid.UUID `json:"attendee_id"`   // FK to Attendee
	ActivityName string    `json:"activity_name"` // FK to Attendee
	ActivityID   uuid.UUID `json:"activity_id"`   // FK to Activity
	ScannedAt    time.Time `json:"scanned_at"`
	Status       string    `json:"status"`     // e.g. "success", "duplicate", "invalid"
	ScannedBy    string    `json:"scanned_by"` // FK to User (staff)
}

type EventInfo struct {
	Event      Event      `json:"event"`
	Activities []Activity `json:"activities"`
}
