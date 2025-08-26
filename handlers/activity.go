package handlers

import (
	"encoding/json"
	"net/http"

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
- 403 Forbidden if user lacks permission
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var c models.ActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	access, err := h.DB.CanCreateActivity(firebaseUser.UID, c.EventID.String())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
	}

	if !access {
		utils.RespondWithError(w, http.StatusForbidden, "You do not have permission to create activity for this event")
		return
	}

	if err := h.DB.CreateActivity(&c); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create Event")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}
