package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
	"github.com/xuri/excelize/v2"
)
/*
Returns:
- 201 Created with created user JSON on success
- 400 Bad Request for invalid input
- 405 Method not allowed except POST
- 409 Failed because User Exists already
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if u.Image_url == "" || u.FullName == "" || u.Position == "" || u.Company == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

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

/*
Returns:
- 200 OK with JSON array of attendees
- 400 Bad Request if event ID is not a valid UUID
- 404 Not Found if event does not exist
- 500 Internal Server Error on database errors
*/
func (h *Handler) GetUsersByEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := mux.Vars(r)["event_id"]
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "event ID not valid")
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
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch attendees")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendees)
}


func (h *Handler) ImportUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streventID := vars["event_id"]
	if streventID == "" {
		http.Error(w, "Missing event_id in URL", http.StatusBadRequest)
		return
	}

	eventID, err := uuid.Parse(streventID)
	if err != nil {
		http.Error(w, "Invalid event_id format", http.StatusBadRequest)
		return
	}

	failedLog := []string{}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Failed to read Excel file", http.StatusBadRequest)
		return
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		http.Error(w, "Failed to read Excel rows", http.StatusInternalServerError)
		return
	}

	var users []models.UserRequest

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) != 4 {
			continue
		}

		if row[0] == "" {
			failedLog = append(failedLog, fmt.Sprintf("Failed to create user in row %v username:%s - company:%s - role:%s - position:%s because a field is empty", i+1, row[0], row[1], row[2], row[3]))
			continue
		}

		user := models.UserRequest{
			Role:     row[0],
			FullName: row[1],
			Position: row[2],
			Company:  row[3],
		}
		users = append(users, user)
	}

	for i, user := range users {
		user, err := h.DB.CreateUser(&models.UserRequest{
			FullName:  user.FullName,
			Company:   user.Company,
			Position:  user.Position,
			Image_url: "https://res.cloudinary.com/dcvr2byrp/image/upload/v1753007426/qocwao1uaykjjnkzqxvo.jpg",
			Role:      user.Role,
			EventId:   eventID.String(),
		})
		if err != nil {
			failedLog = append(failedLog, fmt.Sprintf("Failed to create user in row %v username:%s - company:%s - role:%s - position:%s because %v", i+2, user.FullName, user.Company, user.Role, user.Position, err.Error()))
			continue
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(failedLog)

}
