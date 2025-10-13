package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
)

func (h *Handler) CreateTicketCategory(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var ticketCategoryReq models.TicketCategory
	err := json.NewDecoder(r.Body).Decode(&ticketCategoryReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if ticketCategoryReq.Tag == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ticketCategoryReq.FirebaseID = firebaseUser.UID
	ticketCategoryReq.Type = "tkt"
	ticketCategory, err := h.DB.CreateTicketCategory(ticketCategoryReq)
	if err != nil {
		http.Error(w, "Failed to create ticket category", http.StatusInternalServerError)
		log.Println("Error creating ticket category:", err)
		return
	}
	ticketCategoryRes := models.TicketCategoryResponse{
		ID:          ticketCategory.ID,
		Tag:         ticketCategory.Tag,
		Description: ticketCategory.Description,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ticketCategoryRes)
}

func (h *Handler) CreateInviteeCategory(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var ticketCategoryReq models.TicketCategory
	err := json.NewDecoder(r.Body).Decode(&ticketCategoryReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if ticketCategoryReq.Tag == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	ticketCategoryReq.FirebaseID = firebaseUser.UID
	ticketCategoryReq.Type = "inv"
	ticketCategory, err := h.DB.CreateTicketCategory(ticketCategoryReq)
	if err != nil {
		http.Error(w, "Failed to create ticket category", http.StatusInternalServerError)
		log.Println("Error creating ticket category:", err)
		return
	}
	ticketCategoryRes := models.TicketCategoryResponse{
		ID:          ticketCategory.ID,
		Tag:         ticketCategory.Tag,
		Description: ticketCategory.Description,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticketCategoryRes)
}

func (h *Handler) GetTicketCategories(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	ticketCategories, err := h.DB.GetTicketCategories(firebaseUser.UID, "tkt")
	if err != nil {
		http.Error(w, "Failed to fetch ticket categories", http.StatusInternalServerError)
		log.Println("Error fetching ticket categories:", err)
		return
	}
	ticketCategoriesres := []models.TicketCategoryResponse{}

	for i := range ticketCategories {
		ticketCategoriesres = append(ticketCategoriesres, models.TicketCategoryResponse{
			ID:          ticketCategories[i].ID,
			Tag:         ticketCategories[i].Tag,
			Description: ticketCategories[i].Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticketCategoriesres)
}

func (h *Handler) GetInviteeCategories(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	ticketCategories, err := h.DB.GetTicketCategories(firebaseUser.UID, "inv")
	if err != nil {
		http.Error(w, "Failed to fetch ticket categories", http.StatusInternalServerError)
		log.Println("Error fetching ticket categories:", err)
		return
	}
	ticketCategoriesres := []models.TicketCategoryResponse{}

	for i := range ticketCategories {
		ticketCategoriesres = append(ticketCategoriesres, models.TicketCategoryResponse{
			ID:          ticketCategories[i].ID,
			Tag:         ticketCategories[i].Tag,
			Description: ticketCategories[i].Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticketCategoriesres)
}

// todo: make sure ticket category and event belongs to user creating ticket
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

	if ticketReq.EventID == uuid.Nil.String() || ticketReq.Price <= 0 || ticketReq.Name == "" {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

	ticketCategory, err := h.DB.GetTicketCategory(firebaseUser.UID, ticketReq.TicketCategoryID)
	if err != nil || ticketCategory.Type != "tkt" {
		http.Error(w, "Invalid ticket category id", http.StatusBadRequest)
		return
	}

	ticketReq.FirebaseID = firebaseUser.UID
	ticketReq.Paid = false

	ticket, err := h.DB.CreateTicket(ticketReq)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		log.Println("Error creating ticket:", err)
		return
	}
	ticketRes := models.TicketResponse{
		ID:               ticket.ID,
		TicketCategoryID: ticket.TicketCategoryID,
		EventID:          ticket.EventID,
		Price:            ticket.Price,
		Name:             ticket.Name,
		Paid:             ticket.Paid,
		StartTime:        ticket.StartTime,
		EndTime:          ticket.EndTime,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticketRes)
}

func (h *Handler) CreateInvitee(w http.ResponseWriter, r *http.Request) {
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

	if ticketReq.EventID == uuid.Nil.String() || ticketReq.Name == "" {
		http.Error(w, "Missing or invalid required fields", http.StatusBadRequest)
		return
	}

	ticketCategory, err := h.DB.GetTicketCategory(firebaseUser.UID, ticketReq.TicketCategoryID)
	if err != nil || ticketCategory.Type != "inv" {
		http.Error(w, "Invalid ticket category id", http.StatusBadRequest)
		return
	}

	ticketReq.FirebaseID = firebaseUser.UID
	ticketReq.Price = 0
	ticketReq.Paid = true

	ticket, err := h.DB.CreateTicket(ticketReq)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		log.Println("Error creating ticket:", err)
		return
	}

	inviteeRes := models.InvitationResponse{
		ID:               ticket.ID,
		TicketCategoryID: ticket.TicketCategoryID,
		EventID:          ticket.EventID,
		Name:             ticket.Name,
		StartTime:        ticket.StartTime,
		EndTime:          ticket.EndTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inviteeRes)
}

func (h *Handler) GetTicketsForAllEvents(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	tickets, err := h.DB.GetTicketsForAllEvents(firebaseUser.UID)
	if err != nil {
		http.Error(w, "Failed to fetch tickets", http.StatusInternalServerError)
		log.Println("Error fetching tickets:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

func (h *Handler) GetInviteeForAllEvents(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	invitee, err := h.DB.GetInviteeForAllEvents(firebaseUser.UID)
	if err != nil {
		http.Error(w, "Failed to fetch invitee", http.StatusInternalServerError)
		log.Println("Error fetching invitee:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitee)
}


