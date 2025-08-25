package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

/*
Returns:
- 201 Created with created check-in JSON on success
- 400 Bad Request for invalid input
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var c models.EventRequest
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

func (h *Handler) AddEventWithEventCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: no user in context")
		return
	}

	if len(code) != 6 && len(code) != 7 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid staff_code or admin_code length")
		return
	}

	var event *models.Event
	var err error

	if len(code) == 6 {
		event, err = h.DB.GetEventByStaffId(code)
	} else {
		event, err = h.DB.GetEventByAdminId(code)
	}

	if err != nil {
    log.Println("Failed to fetch event by code:", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch event by code")
		return
	}
	if event == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Event not found")
		return
	}

	err = h.DB.AddStaffToEvent(fireBaseUser.UID, event.ID.String())
	if err != nil {
    log.Println("Failed to fetch event by code:2", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add staff to event")
		return
	}

	json.NewEncoder(w).Encode(event)
}

/*
Returns:
- 200 OK with JSON array of all events
- 500 Internal Server Error if DB query fails
*/
func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	// events, err := h.DB.GetAllEvents()
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		log.Println("Unauthorized: no user in context")
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	events, err := h.DB.GetEventsByFirebaseUser(firebaseUser.UID)
	if err != nil {
		log.Println("Failed to fetch events:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

/*
Returns:
  - 200 OK with JSON { event: {...}, activities: [...] }
  - 400 Bad Request if event_id missing or invalid
  - 404 Not Found if event does not exist
  - 500 Internal Server Error on DB failures
*/

func (h *Handler) GetEventInfo(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

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

	event, err := h.DB.GetEventByFirebaseUser(fireBaseUser.UID, eventID)
	if err != nil {
    log.Println("Failed to fetch event:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch event: "+err.Error())
		return
	}
	if event == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Event not found or no access")
		return
	}

	activities, err := h.DB.GetActivitiesByEvent(fireBaseUser.UID, eventID)
	if err != nil {
    log.Println("Failed to fetch event1 :", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch activities: "+err.Error())
		return
	}

	resp := models.EventInfo{
		Event:      *event,
		Activities: activities,
	}

	log.Print("event info response: ", resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

