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
	"github.com/lib/pq"
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
	// 1️⃣ Get Firebase user from context
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	// 2️⃣ Decode request payload
	var attendeeReq models.AttendeeEnroll
	if err := json.NewDecoder(r.Body).Decode(&attendeeReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 3️⃣ Check if attendee exists for this user + ticket
	attendee, err := h.DB.GetAttendeeFromUserIdandTicketId(attendeeReq.UserID, attendeeReq.TicketID)
	if errors.Is(err, sql.ErrNoRows) {
		// Create new attendee
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
	} else if err != nil {
		log.Println("Error fetching attendee:", err)
		http.Error(w, "Failed to fetch attendee", http.StatusInternalServerError)
		return
	}

	// 4️⃣ Fetch attendee's ticket
	ticket, err := h.DB.GetTicketById(firebaseUser.UID, attendee.TicketID)
	if err != nil {
		log.Println("Error fetching ticket:", err)
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	// 5️⃣ Fetch activity
	activity, err := h.DB.GetActivity(attendeeReq.ActivityID)
	if err != nil {
		log.Println("Error fetching activity:", err)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	// 6️⃣ Ensure activity belongs to the same event as the ticket
	if activity.EventID != uuid.MustParse(ticket.EventID) {
		http.Error(w, "Activity does not belong to the same event as the attendee's ticket", http.StatusBadRequest)
		return
	}

	// 7️⃣ Create attendee_activity row
	attendeeActivityReq := models.AttendeeActivityRequest{
		AttendeeID: uuid.MustParse(attendee.ID),
		ActivityID: attendeeReq.ActivityID,
		FirebaseID: firebaseUser.UID,
	}

	attendeeActivityLog, err := h.DB.CreateAttendeeActivityEnroll(attendeeActivityReq)
	if err != nil {
		// Handle duplicate gracefully (already enrolled)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Println("Attendee already enrolled in activity:", err)
			http.Error(w, "Attendee already enrolled in this activity", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create attendee activity log", http.StatusInternalServerError)
		return
	}

	// 8️⃣ Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendeeActivityLog)
}

func (h *Handler) GetEnrolledAttendee(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	enrolledAttendees, err := h.DB.GetEnrolledAttendee(firebaseUser.UID)
	if err != nil {
		log.Println("Error fetching enrolled attendees:", err)
		http.Error(w, "Failed to fetch enrolled attendees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrolledAttendees)
}

func (h *Handler) GetTicketAttendee(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}
	eventID := r.URL.Query().Get("event_id")

	ticketAttendees, err := h.DB.GetAllTicketAttendee(firebaseUser.UID, eventID)
	if err != nil {
		log.Println("Error fetching ticket attendees:", err)
		http.Error(w, "Failed to fetch ticket attendees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticketAttendees)
}


func (h *Handler) ChangeTicketAttendeeStatus(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	ticketIDStr := r.URL.Query().Get("ticket_id")
	attendeeIDStr := r.URL.Query().Get("attendee_id")
	status := r.URL.Query().Get("status")

	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket_id", http.StatusBadRequest)
		return
	}

	attendeeID, err := uuid.Parse(attendeeIDStr)
	if err != nil {
		http.Error(w, "Invalid attendee_id", http.StatusBadRequest)
		return
	}

	paid := status == "true"

	err = h.DB.ChangeAttendeeTicketStatus(attendeeID, ticketID, paid,firebaseUser.UID)
	if err != nil {
		log.Println("Error changing ticket status:", err)
		http.Error(w, "Failed to update ticket status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Ticket status updated successfully"}`))
}

