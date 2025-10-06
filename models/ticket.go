package models

type Ticket struct {
	ID        string  `json:"id"`
  FirebaseID string  `json:"firebase_id"`
	EventID   string  `json:"event_id"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
	DeletedAt *string `json:"deleted_at"`
}

type TicketRequest struct {
  FirebaseID string  `json:"firebase_id"`
	EventID string  `json:"event_id"`
	Price   float64 `json:"price"`
	Name    string  `json:"name"`
}
