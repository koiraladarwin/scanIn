package models

type Staff struct {
	ID              string `json:"id"`
	FirebaseID      string `json:"firebase_id"`
	StaffGmail      string `json:"staff_gmail"`
	Name            string `json:"name"`
	AutoID          int `json:"auto_id"`
	ImageURL        string `json:"image_url"`
	Phone           string `json:"phone"`
	StaffCategoryID string `json:"staff_category_id"`
	DeletedAt       string `json:"deleted_at"`
	Company         string `json:"company"`
	Position        string `json:"position"`
}

type StaffResponse struct {
	ID              string `json:"id"`
	StaffGmail      string `json:"staff_gmail"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
	Phone           string `json:"phone"`
	StaffCategoryID string `json:"staff_category_id"`
	Company         string `json:"company"`
	Position        string `json:"position"`
}

type StaffCategory struct {
	ID          string `json:"id"`
	FirebaseID  string `json:"firebase_id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	DeletedAt   string `json:"deleted_at"`
}

type StaffCategoryResponse struct {
	ID          string `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	DeletedAt   string `json:"deleted_at"`
}

type StaffEnroll struct {
	ID         string `json:"id"`
	FirebaseID string `json:"firebase_id"`
	StaffID    string `json:"staff_id"`
	EventID    string `json:"event_id"`
	Active     bool   `json:"active"`
	DeletedAt  string `json:"deleted_at"`
}

type EnrolledStaff struct {
	AutoID           int    `json:"auto_id"`
	StaffName        string `json:"staff_name"`
	StaffImage       string `json:"staff_image"`
	EventName        string `json:"event_name"`
	ActivityName     string `json:"activity_name"`
	StafCategoryName string `json:"entry"`
}

type StaffActivities struct {
	ID            string `json:"id"`
	FirebaseID    string `json:"firebase_id"`
	StaffEnrollId string `json:"staff_id"`
	ActivityID    string `json:"activity_id"`
	DeletedAt     string `json:"deleted_at"`
}

type StaffEnrollRequest struct {
	StaffID    string `json:"staff_id"`
	EventID    string `json:"event_id"`
	ActivityID string `json:"session_id"`
}
