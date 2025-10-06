package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

func (h *Handler) CreateUserCategory(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var u models.UsersCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if u.Tag == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	u.FirebaseID = fireBaseUser.UID

	userCat, err := h.DB.CreateUserCategory(&u)

	if errors.Is(err, db.ErrAlreadyExists) {
		utils.RespondWithError(w, http.StatusConflict, "User Category Already Exists")
		return
	}

	if err != nil {
		fmt.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user category")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCat)
}

/*
Returns:
- 201 Created with created user JSON on success
- 400 Bad Request for invalid input
- 405 Method not allowed except POST
- 409 Failed because User Exists already
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var u models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if u.FullName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	u.FirebaseID = fireBaseUser.UID
	user, err := h.DB.CreateUser(&u)

	if errors.Is(err, db.ErrAlreadyExists) {
		utils.RespondWithError(w, http.StatusConflict, "User Already Exists")
		return
	}

	if err != nil {
		fmt.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	users, err := h.DB.GetUsers(fireBaseUser.UID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u models.UserModifyRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if u.FullName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	err := h.DB.UpdateUser(&u)

	if errors.Is(err, db.ErrNotFound) {
		utils.RespondWithError(w, http.StatusNotFound, "User Not Found")
		return
	}

	if err != nil {
		fmt.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

/*
Returns:
- 200 OK with JSON array of attendees
- 400 Bad Request if event ID is not a valid UUID
- 404 Not Found if event does not exist
- 500 Internal Server Error on database errors
*/
func (h *Handler) GetUsersByEvent(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	eventIDStr := mux.Vars(r)["event_id"]
	eventID, err := uuid.Parse(eventIDStr)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "event ID not valid")
		return
	}

	access, err := h.DB.CanSeeAttendee(fireBaseUser.UID, eventIDStr)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check event access")
		return
	}
	if !access {
		utils.RespondWithError(w, http.StatusUnauthorized, "Access denied")
		return
	}

	exists, err := h.DB.EventExists(eventID)
	if err != nil {
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

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}
	userID := r.URL.Query().Get("attendee_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}
	eventID := r.URL.Query().Get("event_id")
	if eventID == "" {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	access, err := h.DB.IsCreator(userID, eventID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check event access")
		return
	}

	if !access {
		utils.RespondWithError(w, http.StatusUnauthorized, "Access denied")
		return
	}

	uuidUser, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteUser(uuidUser)
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	_, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	user, err := h.DB.GetUser(userID)
	if errors.Is(err, db.ErrNotFound) {
		utils.RespondWithError(w, http.StatusNotFound, "User Not Found")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// access, err := h.DB.CanSeeAttendee(fireBaseUser.UID, user.EventId)
	// if err != nil {
	// 	utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check event access")
	// 	return
	// }
	// if !access {
	// 	utils.RespondWithError(w, http.StatusUnauthorized, "Access denied")
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
