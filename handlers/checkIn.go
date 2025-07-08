package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

/*
CreateCheckIn accepts JSON:
{
  "attendee_id": "uuid-string",
  "activity_id": "uuid-string",
  "scanned_at": "timestamp",
  "status": "string",
  "scanned_by": "string"
}

Returns:
- 201 Created with created check-in JSON on success
- 400 Bad Request for invalid input
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateCheckIn(w http.ResponseWriter, r *http.Request) {
	var c models.CheckInLog
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := h.DB.CreateCheckInLog(&c); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check in")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

/*
CreateCheckIn2 checks if attendee already checked in before creating.

Accepts same JSON as CreateCheckIn.

Returns:
- 201 Created with created check-in JSON on success
- 400 Bad Request for invalid input
- 409 Conflict if already checked in for that attendee and activity
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateCheckIn2(w http.ResponseWriter, r *http.Request) {
	var c models.CheckInLog
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	exists, err := h.DB.CheckInExists(c.AttendeeID, c.ActivityID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check existing check-in")
		return
	}
	if exists {
		utils.RespondWithError(w, http.StatusConflict, "Already checked in for this activity")
		return
	}

	if err := h.DB.CreateCheckInLog(&c); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check in")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

/*
ModifyCheckIn updates an existing check-in by ID.

Accepts JSON:
{
  "attendee_id": "uuid-string",
  "activity_id": "uuid-string",
  "scanned_at": "timestamp",
  "status": "string",
  "scanned_by": "string"
}

Returns:
- 200 OK with updated check-in JSON on success
- 400 Bad Request for invalid ID or input
- 500 Internal Server Error on DB failure
*/
func (h *Handler) ModifyCheckIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var c models.CheckInLog
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	c.ID = id

	if err := h.DB.UpdateCheckInLog(&c); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update")
		return
	}

	json.NewEncoder(w).Encode(c)
}

