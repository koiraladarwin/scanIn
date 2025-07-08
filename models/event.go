package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered person in the system
type User struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
}

type UserWithRole struct {
	ID       uuid.UUID `json:"id"`
	Role     string `json:"role"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
}

// Event represents a specific event like "Teej 2025"
type Event struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
}

// Activity is a feature or sub-event within an event (e.g. Registration, Food)
type Activity struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"` // FK to Event
	Name      string    `json:"name"`     // e.g. Registration, Dinner
	Type      string    `json:"type"`     // e.g. "registration", "food", etc.
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// Attendee links a user to an event
type Attendee struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`  // FK to User
	EventID uuid.UUID `json:"event_id"` // FK to Event
	Role    string    `json:"role"`     // "participant", "staff", "member"
}

// CheckInLog is created when an attendee scans at an activity
type CheckInLog struct {
	ID         uuid.UUID `json:"id"`
	AttendeeID uuid.UUID `json:"attendee_id"` // FK to Attendee
	ActivityID uuid.UUID `json:"activity_id"` // FK to Activity
	ScannedAt  time.Time `json:"scanned_at"`
	Status     string    `json:"status"`     // e.g. "success", "duplicate", "invalid"
	ScannedBy  uuid.UUID `json:"scanned_by"` // FK to User (staff)
}

type EventWithActivities struct {
	Event      Event      `json:"event"`
	Activities []Activity `json:"activities"`
}
