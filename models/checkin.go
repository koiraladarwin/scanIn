package models

import (
	"github.com/google/uuid"
	"time"
)

type CheckInLog struct {
	ID         uuid.UUID  `json:"id"`
	AttendeeId uuid.UUID  `json:"user_id"`
	ActivityID uuid.UUID  `json:"activity_id"`
	ScannedAt  time.Time  `json:"scanned_at"`
	ScannedBy  string     `json:"scanned_by"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type CheckInRespose struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"full_name"`
	AutoId       string    `json:"auto_id"`
	Role         string    `json:"role"`
	AttendeeId   uuid.UUID `json:"user_id"`
	ActivityName string    `json:"activity_name"`
	ActivityID   uuid.UUID `json:"activity_id"`
	ScannedAt    time.Time `json:"scanned_at"`
	ScannedBy    string    `json:"scanned_by"`
}

type CheckInLogRequest struct {
	AttendeeId uuid.UUID `json:"attendee_id"`
	ActivityID uuid.UUID `json:"activity_id"`
}
