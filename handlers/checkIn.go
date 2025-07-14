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
CreateCheckIn accepts JSON:

	{
	  "attendee_id": "uuid-string",
	  "activity_id": "uuid-string",
	  "scanned_at": "timestamp",
	  "status": "checked",
	  "scanned_by": "string"
	}

Returns:
- 201 Created with created check-in JSON on success
- 400 Bad Request for invalid input
- 409 Conflict if already checked in for that attendee and activity
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateCheckIn(w http.ResponseWriter, r *http.Request) {
	var c models.CheckInLog
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.DB.CheckInExists(c.AttendeeID, c.ActivityID)
	if errors.Is(err, db.ErrNotFound) {
		err := h.DB.CreateCheckInLog(&c)
		if err != nil {
			log.Print(err.Error())
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check in")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check existing check-in")
		return
	}

	user, err := h.DB.GetCheckInLog(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch existing check-in")
		return
	}
	if user.Status == "checked" {
		utils.RespondWithError(w, http.StatusConflict, "Cannot Check in twice")
		return
	}

	user.Status = "checked"
	err = h.DB.UpdateCheckInLog(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

/*
GetCheckIn , retrives all check Ins

Returns:
- 200 OK with updated check-in JSON on success
- 500 Internal Server Error on DB failure
*/

func (h *Handler) GetCheckIn(w http.ResponseWriter, r *http.Request) {

	checkInLogs, err := h.DB.GetAllCheckInLog()
	if err != nil {
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "cant get checkinlogs")
		return
	}

	user, err := h.DB.GetUserByAttendeeid(checkInLogs.AttendeeID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cant get user details")
		return
	}

	scannedBy, err := h.DB.GetUserByAttendeeid(checkInLogs.ScannedBy)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cant get user details")
		return
	}

	checkInResponse := models.CheckInRespose{
		ID:         checkInLogs.ID,
		FullName:   user.FullName,
		AttendeeID: checkInLogs.AttendeeID,
		ActivityID: checkInLogs.ActivityID,
		ScannedAt:  checkInLogs.ScannedAt,
		ScannedBy:  checkInLogs.ScannedBy,
		Status:     checkInLogs.Status,
    ScannedByName: scannedBy.FullName,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkInResponse)
}
