package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
)

// todo: make sure the userid and eventid belong to the same firebase user
func (h *Handler) CreateAttendee(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
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
	attendeeReq.FirebaseID = firebaseUser.UID
	if attendeeReq.UserID == uuid.Nil || attendeeReq.TicketID == uuid.Nil {
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
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var attendeeActivityReq models.AttendeeActivityRequest
	err := json.NewDecoder(r.Body).Decode(&attendeeActivityReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if attendeeActivityReq.AttendeeID == uuid.Nil || attendeeActivityReq.ActivityID == uuid.Nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	attendeeActivityReq.FirebaseID = firebaseUser.UID
	activity, err := h.DB.GetActivity(attendeeActivityReq.ActivityID)
	if err != nil {
		log.Println("Error fetching activity:", err.Error())
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}
	event, err := h.DB.GetEventFromAttendee(attendeeActivityReq.AttendeeID)
	if err != nil {
		http.Error(w, "Event not found for attendee", http.StatusNotFound)
		return
	}

	if activity.EventID != event.ID {
		http.Error(w, "Activity does not belong to the same event as the attendee", http.StatusBadRequest)
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

func (h *Handler) EnrollAttendee(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var attendeeReq models.AttendeeEnroll
	err := json.NewDecoder(r.Body).Decode(&attendeeReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	attendee, err := h.DB.GetAttendeeFromUserIdandEventId(attendeeReq.UserID, attendeeReq.EventID)
	if errors.Is(err, sql.ErrNoRows) {
		newAttendeeReq := models.AttendeeRequest{
			UserID:     attendeeReq.UserID,
			TicketID:   attendeeReq.TicketID,
			FirebaseID: firebaseUser.UID,
		}
		attendee, err = h.DB.CreateAttendee(newAttendeeReq)
		if err != nil {
			log.Println("Error creating attendee:", err)
			http.Error(w, "Failed to create attendee", http.StatusInternalServerError)
			return
		}
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error fetching attendee:", err)
		http.Error(w, "Failed to fetch attendee", http.StatusInternalServerError)
		return
	}
	activity, err := h.DB.GetActivity(attendeeReq.ActivityID)
	if err != nil {
		log.Println("Error fetching activity:", err.Error())
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}
	if activity.EventID != attendeeReq.EventID {
		http.Error(w, "Activity does not belong to the same event as the attendee", http.StatusBadRequest)
		return

	}
	ticket, err := h.DB.GetTicketById(firebaseUser.UID, attendee.TicketID)
	if err != nil {
		log.Println("Error fetching ticket:", err.Error())
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}
	if ticket.EventID != attendeeReq.EventID.String() {
		http.Error(w, "Ticket does not belong to the same event as the attendee", http.StatusBadRequest)
		return
	}
	attendeeActivityReq := models.AttendeeActivityRequest{
		AttendeeID: uuid.MustParse(attendee.ID),
		ActivityID: attendeeReq.ActivityID,
		FirebaseID: firebaseUser.UID,
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
