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

type Attendee struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	EventID uuid.UUID `json:"event_id"`
}

type Event struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Location            string    `json:"location"`
	NumberOfParticipant int       `json:"number_of_participant"`
	NumberOfStaff       int       `json:"number_of_staff"`
	StaffCode           *string   `json:"staff_code"`
}

type EventCreateRequest struct {
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Location            string    `json:"location"`
	NumberOfParticipant int       `json:"number_of_participant"`
}

type EventModifyRequest struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Location            string    `json:"location"`
}

type EventInfo struct {
	Event      Event      `json:"event"`
	Activities []Activity `json:"activities"`
}
