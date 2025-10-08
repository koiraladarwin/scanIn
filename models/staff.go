package models

type Staff struct {
	ID              string `json:"id"`
	FirebaseID      string `json:"firebase_id"`
	StaffGmail      string `json:"staff_gmail"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
	Phone           string `json:"phone"`
	StaffCategoryID string `json:"staff_category_id"`
	DeletedAt       string `json:"deleted_at"`
}

type StaffResponse struct {
	ID              string `json:"id"`
	StaffGmail      string `json:"staff_gmail"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
	Phone           string `json:"phone"`
	StaffCategoryID string `json:"staff_category_id"`
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

type StaffActivities struct {
  ID         string `json:"id"`
  FirebaseID string `json:"firebase_id"`
  StaffID    string `json:"staff_id"`
  ActivityID string `json:"activity_id"`
  Role       string `json:"role"`
  DeletedAt  string `json:"deleted_at"`
}
