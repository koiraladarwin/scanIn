package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

func (h *Handler) CreateStaffCategory(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	var u models.StaffCategory
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if u.Tag == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	u.FirebaseID = fireBaseUser.UID

	userCat, err := h.DB.CreateStaffCategory(&u)

	if errors.Is(err, db.ErrAlreadyExists) {
		utils.RespondWithError(w, http.StatusConflict, "User Category Already Exists")
		return
	}

	if err != nil {
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user category")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCat)
}

func (h *Handler) GetStaffCategories(w http.ResponseWriter, r *http.Request) {
	fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	userCategories, err := h.DB.GetStaffCategories(fireBaseUser.UID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user categories")
		return
	}
	userCategoriesResponse := make([]models.StaffCategoryResponse, len(userCategories))

	for i, uc := range userCategories {
		userCategoriesResponse[i] = models.StaffCategoryResponse{
			ID:          uc.ID,
			Tag:         uc.Tag,
			Description: uc.Description,
			DeletedAt:   uc.DeletedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userCategoriesResponse)
}

func (h *Handler) CreateStaff(w http.ResponseWriter, r *http.Request) {
  fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
  if !ok {
    http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
    return
  }

  var u models.Staff
  if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
    utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
    return
  }

  if u.Name == "" || u.StaffGmail == "" || u.Phone == "" || u.StaffCategoryID == "" {
    utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
    return
  }

  u.FirebaseID = fireBaseUser.UID

  staff, err := h.DB.CreateStaff(&u)

  if errors.Is(err, db.ErrAlreadyExists) {
    utils.RespondWithError(w, http.StatusConflict, "Staff Already Exists")
    return
  }

  if err != nil {
    log.Print(err.Error())
    utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create staff")
    return
  }

  staffRes := models.StaffResponse{
    ID:              staff.ID,
    StaffGmail:      staff.StaffGmail,
    Name:            staff.Name,
    ImageURL:        staff.ImageURL,
    Phone:           staff.Phone,
    StaffCategoryID: staff.StaffCategoryID,
  }

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(staffRes)

}

func (h *Handler) GetStaffs(w http.ResponseWriter, r *http.Request) {
  fireBaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
  if !ok {
    http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
    return
  }

  staffs, err := h.DB.GetStaffs(fireBaseUser.UID)
  if err != nil {
    utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch staffs")
    return
  }
  staffsResponse := make([]models.StaffResponse, len(staffs))
  for i := range staffs { 
    staffsResponse[i] = models.StaffResponse{
      ID:              staffs[i].ID,
      StaffGmail:      staffs[i].StaffGmail,
      Name:            staffs[i].Name,
      ImageURL:        staffs[i].ImageURL,
      Phone:           staffs[i].Phone,
      StaffCategoryID: staffs[i].StaffCategoryID,
    }
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(staffsResponse)
}

