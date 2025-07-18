package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

/*
RegisterAttendee accepts JSON:

	{
	  "user_id": "uuid-string",
	  "event_id": "uuid-string",
		"role": "participant|staff|member"
	}

Returns:
- 201 Created with attendee JSON on success
- 400 Bad Request for invalid JSON or missing fields
- 409 Conflict if attendee already registered
- 500 Internal Server Error on other DB failures
*/
func (h *Handler) RegisterAttendee(w http.ResponseWriter, r *http.Request) {
	var a models.Attendee

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if a.UserID == uuid.Nil || a.EventID == uuid.Nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user_id or event_id")
		return
	}

	err := h.DB.CreateAttendee(&a)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			utils.RespondWithError(w, http.StatusConflict, "Attendee already registered")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to register attendee")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

/*
GetUsersByEvent
Fetches list of attendees for given event ID (UUID in URL path param).

Returns:
- 200 OK with JSON array of attendees
- 400 Bad Request if event ID is not a valid UUID
- 404 Not Found if event does not exist
- 500 Internal Server Error on database errors
*/
func (h *Handler) GetUsersByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := mux.Vars(r)["event_id"]
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "event ID not valid")
		return
	}

	exists, err := h.DB.EventExists(eventID)
	if err != nil {
    log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "database error")
		return
	}

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "event not found")
		return
	}

	attendees, err := h.DB.GetUsersByEvent(eventID)
	if err != nil {
    log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch attendees")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendees)
}
