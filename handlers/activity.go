package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	"github.com/koiraladarwin/scanin/features/firebaseauth"
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
	//todo
	//change this route so only super admin and event creator can create activity
	firebaseId, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if !slices.Contains(superAdminsEmails, firebaseId.Email) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var c models.ActivityCreateRequest
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
UpdateActivity accepts JSON:

	{
	  "name": "string",
	  "type": "string",
	  "start_time": "2025-07-08T15:30:00Z",
	  "end_time": "2025-07-09T15:30:00Z",
	  "location": "string"
	}

Path Param:

	activity_id (uuid-string)

Returns:
- 200 OK with updated activity JSON on success
- 400 Bad Request for invalid input
- 404 Not Found if activity doesn’t exist
- 500 Internal Server Error on DB failure
*/
func (h *Handler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	var activity models.Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	err = h.DB.UpdateActivity(&activity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "Activity not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update activity")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activity)
}
