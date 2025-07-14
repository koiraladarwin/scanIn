package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

/*
CreateEvent accepts JSON:
{
  "name": "string",
  "description": "timestamp",
  "start_time": "2025-07-08T15:30:00Z",
  "end_time": "2025-07-09T15:30:00Z",
  "location": "string"
}

Returns:
- 201 Created with created check-in JSON on success
- 400 Bad Request for invalid input
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var c models.Event
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := h.DB.CreateEvent(&c); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create Event")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

/*
GetEvent handles GET 

Returns:
- 200 OK with JSON array of all events
- 500 Internal Server Error if DB query fails
*/
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	events, err := h.DB.GetAllEvents()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}


