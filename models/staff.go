package models

type Staff struct {
	FireBaseId     string `json:"firebase_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ImageUrl       string `json:"image_url"`
	CanSeeScanned  bool   `json:"can_see_scanned"`
	CanCreateAttendee bool   `json:"can_add_attendee"`
	CanSeeAttendee bool   `json:"can_see_attendee"`
}
