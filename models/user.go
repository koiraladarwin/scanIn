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

type User struct {
	ID              uuid.UUID `json:"id"`
	FirebaseID      string    `json:"firebase_id"`
	FullName        string    `json:"full_name"`
	Company         string    `json:"company"`
	Position        string    `json:"position"`
	Image_url       string    `json:"image_url"`
	AutoId          int       `json:"auto_id"`
	UsersCategoryID string    `json:"users_category_id"`
}

type UserModifyRequest struct {
	ID              uuid.UUID `json:"id"`
	FullName        string    `json:"full_name"`
	Company         string    `json:"company"`
	Position        string    `json:"position"`
	Image_url       string    `json:"image_url"`
	AutoId          int       `json:"auto_id"`
	UsersCategoryID string    `json:"users_category_id"`
}

type UserRequest struct {
	FirebaseID      string `json:"firebase_id"`
	FullName        string `json:"full_name"`
	Company         string `json:"company"`
	Position        string `json:"position"`
	Image_url       string `json:"image_url"`
	UsersCategoryID string `json:"users_category_id"`
}
