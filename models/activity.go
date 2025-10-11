package models

import (
	"github.com/google/uuid"
	"time"
)

type Activity struct {
	ID                  uuid.UUID `json:"id"`
	EventID             uuid.UUID `json:"event_id"`
	FirebaseID          string    `json:"firebase_id"`
	Name                string    `json:"name"`
	HallName            string    `json:"hall_name"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	NumberOfScanedUsers int       `json:"number_of_scaned_users"`
}
type ActivityDetails struct {
	ID              uuid.UUID `json:"id"`
	EventID         uuid.UUID `json:"event_id"`
	Name            string    `json:"name"`
	HallName        string    `json:"hall_name"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	InvitationCount int       `json:"invitations_count"`
	TicketCount     int       `json:"ticket_count"`
}

type ActivityCreateRequest struct {
	EventID    uuid.UUID `json:"event_id"`
	FirebaseID string    `json:"firebase_id"`
	Name       string    `json:"name"`
	HallName   string    `json:"hall_name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}
