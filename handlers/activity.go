package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

/*
CreateActivity accepts JSON:

	{
	  "event_id": "uuid-string",
	  "name": "string",
	  "type": "string",
	  "start_time": "2025-07-08T15:30:00Z",
	  "end_time": "2025-07-09T15:30:00Z",
	  "location": "string"
	}

Returns:
- 201 Created with created check-in JSON on success
- 400 Bad Request for invalid input
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var c models.Activity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := h.DB.CreateActivity(&c); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create Event")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

/*
GetEventInfo returns event details along with its activities.
Accepts:
  - event_id (query param): UUID of the event

Returns:
  - 200 OK with JSON { event: {...}, activities: [...] }
  - 400 Bad Request if event_id missing or invalid
  - 404 Not Found if event does not exist
  - 500 Internal Server Error on DB failures
*/
func (h *Handler) GetEventInfo(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing event_id query parameter")
		return
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid event_id format")
		return
	}

	event, err := h.DB.GetEvent(eventID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch event"+err.Error())
		return
	}

	if event == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Event not found")
		return
	}

	numberOfParticipants, err := h.DB.GetNumberOfAttendeesByEvent(eventID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch number of participants")
		return
	}

	event.NumberOfParticipant = numberOfParticipants

	activities, err := h.DB.GetActivitiesByEvent(eventID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch activities"+err.Error())
		return
	}

	resp := models.EventInfo{
		Event:      *event,
		Activities: activities,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
