package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

func (h *Handler) GiveRoleToStaff(w http.ResponseWriter, r *http.Request) {
	editRoleReq := &models.RoleRequest{}
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: no user in context")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(editRoleReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	isCreator, err := h.DB.IsCreator(fireBaseUser.UID, editRoleReq.EventId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check creator status")
		return
	}

	if !isCreator {
		utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to give roles")
		return
	}

	if err := h.DB.AddEventRole(*editRoleReq); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create role")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ModifyRoleToStaff(w http.ResponseWriter, r *http.Request) {
	createRoleReq := &models.EditRoleRequest{}
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: no user in context")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(createRoleReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	isCreator, err := h.DB.IsCreator(fireBaseUser.UID, createRoleReq.EventId)

	if err != nil {
    log.Printf("Failed to check creator status for user %s in event %s: %v", fireBaseUser.UID, createRoleReq.EventId, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check creator status")
		return
	}

	if !isCreator {
		utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to give roles")
		return
	}

	if err := h.DB.ModifyEventRole(*createRoleReq); err != nil {
    log.Printf("Failed to modify role for user %s in event %s: %v", createRoleReq.FireBaseId, createRoleReq.EventId, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create role")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetStaffsByEvent(w http.ResponseWriter, r *http.Request) {
	eventId := mux.Vars(r)["event_id"]
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
  staffs := []models.Staff{}
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: no user in context")
		return
	}

	isCreator, err := h.DB.IsCreator(fireBaseUser.UID, eventId)

	if err != nil {
    log.Printf("Failed to check creator status for user %s in event %s: %v", fireBaseUser.UID, eventId, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check creator status")
		return
	}

	if !isCreator {
		utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to give roles")
		return
	}

	staffsFirebaseIds, err := h.DB.GetStaffByEvent(eventId)
	if err != nil {
    log.Printf("Failed to get staffs for event %s: %v", eventId, err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get staffs")
		return
	}

  ctx := r.Context()
  for _, staffId := range staffsFirebaseIds {
    staff  := models.Staff{}

    firebaseUser, err := h.FbAuth.GetUserByIdToken(ctx,staffId.FireBaseId)
    if err != nil { 
      log.Printf("Failed to fetch Firebase user info for ID %s: %v", staffId.FireBaseId, err)
      continue 
    }

    staff.Name = firebaseUser.DisplayName
    staff.FireBaseId = firebaseUser.UID
    staff.ImageUrl = firebaseUser.PhotoURL 
    staff.Email = firebaseUser.Email
    staff.CanSeeScanned = staffId.CanSeeScanned
    staff.CanCreateAttendee = staffId.CanCreateAttendee
    staff.CanSeeAttendee = staffId.CanSeeAttendee

    staffs = append(staffs, staff) 
  }

	w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(staffs)
}
