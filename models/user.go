package models

import "github.com/google/uuid"

type UsersCategory struct {
	ID          uuid.UUID `json:"id"`
	FirebaseID  string    `json:"firebase_id"`
	Tag         string    `json:"tag"`
	Description string    `json:"description"`
	DeletedAt   string    `json:"deleted_at"`
}

type UsersCategoryRequest struct {
	FirebaseID  string `json:"firebase_id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

type UsersCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Tag         string    `json:"tag"`
	Description string    `json:"description"`
	DeletedAt   string    `json:"deleted_at"`
}

type User struct {
	ID              uuid.UUID `json:"id"`
	FirebaseID      string    `json:"firebase_id"`
	FullName        string    `json:"full_name"`
	Company         string    `json:"company"`
	Position        string    `json:"position"`
	Image_url       string    `json:"image_url"`
	AutoId          int       `json:"auto_id"`
	PhoneNumber     string    `json:"phone_number"`
	UsersCategoryID string    `json:"attendee_category_id"`
}

type UserModifyRequest struct {
	ID              uuid.UUID `json:"id"`
	FullName        string    `json:"full_name"`
	Company         string    `json:"company"`
	Position        string    `json:"position"`
	Image_url       string    `json:"image_url"`
	AutoId          int       `json:"auto_id"`
	PhoneNumber     string    `json:"phone_number"`
	UsersCategoryID string    `json:"attendee_category_id"`
}

type UserRequest struct {
	FirebaseID      string `json:"firebase_id"`
	FullName        string `json:"full_name"`
	Company         string `json:"company"`
	Position        string `json:"position"`
	Image_url       string `json:"image_url"`
	PhoneNumber     string `json:"phone_number"`
	UsersCategoryID string `json:"attendee_category_id"`
}
