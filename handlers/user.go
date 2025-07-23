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
CreateUser accepts JSON:

	{
	  "full_name": "string",
	  "email": "string",
	  "phone": "string",
	}

Returns:
- 201 Created with created user JSON on success
- 400 Bad Request for invalid input
- 405 Method not allowed except POST
- 409 Failed because User Exists already
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
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
GetAllUsers returns a list of all users.

Returns:
- 200 OK with JSON array of users on success
- 500 Internal Server Error on DB failure
*/
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.DB.GetAllUsers()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

/*
CreateUser accepts JSON:

	{
	  "full_name": "string",
	  "company": "string",
	  "position": "string",
    "image_url":"string",
    "role": "string"
	}

Returns:
- 201 Created with created user JSON on success
- 400 Bad Request for invalid input
- 405 Method not allowed except POST
- 409 Failed because User Exists already
- 500 Internal Server Error on DB failure
*/

// CreateUser2 is a modified version of CreateUser that accepts an event ID from the URL path.
func (h *Handler) CreateUser2(w http.ResponseWriter, r *http.Request) {
	var userReq models.UserRequest
	vars := mux.Vars(r)
	idStr := vars["event_id"]
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Event ID")
		return
	}
	EventUuid, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Event ID")
	}

	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if userReq.Image_url == "" || userReq.FullName == "" || userReq.Position == "" || userReq.Company == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := h.DB.CreateUser(&models.User{
		FullName:  userReq.FullName,
		Company:   userReq.Company,
		Position:  userReq.Position,
		Image_url: userReq.Image_url,
		Role:      userReq.Role,
	})

	if errors.Is(err, db.ErrAlreadyExists) {
		utils.RespondWithError(w, http.StatusConflict, "User Already Exists")
		return
	}

	if err != nil {
		fmt.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	attendde, err := h.DB.CreateAttendee(&models.Attendee{
		UserID:  user.ID,
		EventID: EventUuid,
	})

	user.ID = attendde.ID

	if err != nil {
		fmt.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "User created by failed to create Attendee . Check the event Id dude")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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

		if row[0] == "" || row[1] == "" || row[2] == "" || row[3] == "" {
			failedLog = append(failedLog, fmt.Sprintf("Failed to create user in row %v username:%s - company:%s - role:%s - position:%s because a field is empty", i+1, row[0], row[1], row[2], row[3]))
			continue
		}

		user := models.UserRequest{
			FullName: row[0],
			Company:  row[1],
			Role:     row[2],
			Position: row[3],
		}
		users = append(users, user)
	}

	for i, user := range users {
		user, err := h.DB.CreateUser(&models.User{
			FullName:  user.FullName,
			Company:   user.Company,
			Position:  user.Position,
			Image_url: "https://res.cloudinary.com/dcvr2byrp/image/upload/v1753007426/qocwao1uaykjjnkzqxvo.jpg",
			Role:      user.Role,
		})
		if err != nil {
			failedLog = append(failedLog, fmt.Sprintf("Failed to create user in row %v username:%s - company:%s - role:%s - position:%s because %v", i+2, user.FullName, user.Company, user.Role, user.Position, err.Error()))
			continue
		}
		_, err = h.DB.CreateAttendee(&models.Attendee{
			EventID: eventID,
			UserID:  user.ID,
		})

		if err != nil {
			failedLog = append(failedLog, fmt.Sprintf("failed to add addendee please check the event id"))
			continue
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(failedLog)

}
