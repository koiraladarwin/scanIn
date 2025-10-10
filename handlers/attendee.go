package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
)

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

func (h *Handler) CreateAttendeeActivityEnroll(w http.ResponseWriter, r *http.Request) {
  _, ok := firebaseauth.FbUserFromContext(r.Context())
  if !ok {
    http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
    return
  }

  var attendeeActivityReq models.AttendeeActivity
  err := json.NewDecoder(r.Body).Decode(&attendeeActivityReq)
  if err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }

  if  attendeeActivityReq.AttendeeID == uuid.Nil || attendeeActivityReq.ActivityID == uuid.Nil {
    http.Error(w, "Missing required fields", http.StatusBadRequest)
    return
  }

  attendeeActivityLog, err := h.DB.CreateAttendeeActivityEnroll(attendeeActivityReq)
  if err != nil {
    log.Println("Error creating attendee activity log:", err)
    http.Error(w, "Failed to create attendee activity log", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(attendeeActivityLog)

}
