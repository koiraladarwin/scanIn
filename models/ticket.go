package models

import "github.com/google/uuid"

type TicketCategory struct {
	ID          string  `json:"id"`
	FirebaseID  string  `json:"firebase_id"`
	Tag         string  `json:"tag"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	DeletedAt   *string `json:"deleted_at"`
}

type TicketCategoryResponse struct {
	ID          string `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

type Ticket struct {
	ID               string  `json:"id"`
	FirebaseID       string  `json:"firebase_id"`
	TicketCategoryID string  `json:"ticket_category_id"`
	EventID          string  `json:"event_id"`
	Price            float64 `json:"price"`
	Paid             bool    `json:"paid"`
	Name             string  `json:"name"`
	DeletedAt        *string `json:"deleted_at"`
}

type TicketRequest struct {
	FirebaseID       string    `json:"firebase_id"`
	TicketCategoryID uuid.UUID `json:"ticket_category_id"`
	EventID          string    `json:"event_id"`
	Price            float64   `json:"price"`
	Name             string    `json:"name"`
	Paid             bool      `json:"paid"`
}

type TicketResponse struct {
	ID               string  `json:"id"`
	TicketCategoryID string  `json:"ticket_category_id"`
	EventID          string  `json:"event_id"`
	Price            float64 `json:"price"`
	Name             string  `json:"name"`
	Paid             bool    `json:"paid"`
}

type InviteeResponse struct {
	ID               string `json:"id"`
	TicketCategoryID string `json:"ticket_category_id"`
	EventID          string `json:"event_id"`
	Name             string `json:"name"`
}
