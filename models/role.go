package models

type RoleRequest struct {
	EventId        string `json:"event_id"`
	FireBaseId     string `json:"firebase_id"`
	CanAddAttendee bool   `json:"can_add_attendee"`
	CanSeeAttendee bool   `json:"can_see_attendee"`
	CanSeeScanned  bool   `json:"can_see_scanned"`
}

type EditRoleRequest struct {
	FireBaseId     string `json:"firebase_id"`
	EventId        string `json:"event_id"`
	CanAddAttendee bool   `json:"can_add_attendee"`
	CanSeeAttendee bool   `json:"can_see_attendee"`
	CanSeeScanned  bool   `json:"can_see_scanned"`
}

type Role struct {
	ID                string `json:"id"`
	IsCreator         bool   `json:"is_creator"`
	CanSeeScanned     bool   `json:"can_see_scanned"`
	CanCreateActivity bool   `json:"can_create_activity"`
	CanCreateAttendee bool   `json:"can_create_attendee"`
	CanSeeAttendee    bool   `json:"can_see_attendee"`
}
