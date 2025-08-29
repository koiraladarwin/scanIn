package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Company   string    `json:"company"`
	Position  string    `json:"position"`
	Image_url string    `json:"image_url"`
	AutoId    int       `json:"auto_id"`
	EventId   string    `json:"event_id"`
	Role      string    `json:"role"`
}

type UserModifyRequest struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Company   string    `json:"company"`
	Position  string    `json:"position"`
	Image_url string    `json:"image_url"`
	AutoId    int       `json:"auto_id"`
	EventId   string    `json:"event_id"`
	Role      string    `json:"role"`
}

type UserRequest struct {
	FullName  string `json:"full_name"`
	Company   string `json:"company"`
	Position  string `json:"position"`
	Image_url string `json:"image_url"`
	EventId   string `json:"event_id"`
	Role      string `json:"role"`
}
