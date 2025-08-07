package models

type Role struct {
	ID                string `json:"id"`
	IsCreator         bool   `json:"is_creator"`
	CanSeeScanned     bool   `json:"can_see_scanned"`
	CanCreateActivity bool   `json:"can_create_activity"`
	CanCreateAttendee bool   `json:"can_create_attendee"`
	CanSeeAttendee    bool   `json:"can_see_attendee"`
}
