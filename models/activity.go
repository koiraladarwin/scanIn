package models

import (
	"time"
	"github.com/google/uuid"
)

type Activity struct {
	ID                  uuid.UUID `json:"id"`
	EventID             uuid.UUID `json:"event_id"`
	Name                string    `json:"name"`     
	Type                string    `json:"type"`    
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	NumberOfScanedUsers int       `json:"number_of_scaned_users"`
}

type ActivityCreateRequest struct {
	EventID             uuid.UUID `json:"event_id"`
	Name                string    `json:"name"`     
	Type                string    `json:"type"`    
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	NumberOfScanedUsers int       `json:"number_of_scaned_users"`
}


