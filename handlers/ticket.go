package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
)

func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var ticketReq models.TicketRequest
	err := json.NewDecoder(r.Body).Decode(&ticketReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if ticketReq.EventID == uuid.Nil.String() || ticketReq.Price < 0 || ticketReq.Name == "" {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

  ticketReq.FirebaseID = firebaseUser.UID

	ticket, err := h.DB.CreateTicket(ticketReq)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
    log.Println("Error creating ticket:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}
