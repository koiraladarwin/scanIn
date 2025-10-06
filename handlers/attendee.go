package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
)

func (h *Handler) CreateAttendeeCategory(w http.ResponseWriter, r *http.Request){
  firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context()) 
  if !ok {
    http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
    return
  }

  var attendeeCategoryReq models.AttendeeCategory
  err := json.NewDecoder(r.Body).Decode(&attendeeCategoryReq)
  if err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }
  if attendeeCategoryReq.Tag == "" || attendeeCategoryReq.Description == "" {
    http.Error(w, "Missing required fields", http.StatusBadRequest)
    return
  }
  attendeeCategoryReq.FirebaseID = firebaseUser.UID
  attendeeCategory, err := h.DB.CreateAttendeeCategory(attendeeCategoryReq)
  if err != nil {
    http.Error(w, "Failed to create attendee category", http.StatusInternalServerError)
    log.Println("Error creating attendee category:", err)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(attendeeCategory)

}

func (h *Handler) GetAttendeeCategories(w http.ResponseWriter, r *http.Request){
  firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context()) 
  if !ok {
    http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
    return
  }
  categories, err := h.DB.GetAttendeeCategories(firebaseUser.UID)
  if err != nil {
    http.Error(w, "Failed to get attendee categories", http.StatusInternalServerError)
    log.Println("Error getting attendee categories:", err)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(categories) 
}

//todo: make sure the userid and eventid belong to the same firebase user
func (h *Handler) CreateAttendee(w http.ResponseWriter, r *http.Request) {
	_, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var attendeeReq models.AttendeeRequest
	err := json.NewDecoder(r.Body).Decode(&attendeeReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if  attendeeReq.UserID == uuid.Nil || attendeeReq.TicketID == uuid.Nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	attendee, err := h.DB.CreateAttendee(attendeeReq)
	if err != nil {
    log.Println("Error creating attendee:", err)
		http.Error(w, "Failed to create attendee", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendee)

}
