package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
	"github.com/xuri/excelize/v2"
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
	userVal := r.Context().Value(firebaseauth.FirebaseUserContextKey)
	if userVal == nil {
		http.Error(w, "Unauthorized: no user in context", http.StatusUnauthorized)
		return
	}

	fbuser, ok := userVal.(*auth.UserRecord)
	if !ok {
		http.Error(w, "Context value is not of type *auth.UserRecord", http.StatusInternalServerError)
		return
	}
	var c models.CheckInLogRequest

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.DB.CheckInExists(c.UserID, c.ActivityID)
	if errors.Is(err, db.ErrNotFound) {
		scannedBy := fbuser.Email
		scannedAt := time.Now()
		err := h.DB.CreateCheckInLog(&models.CheckInLog{
			UserID:     c.UserID,
			ActivityID: c.ActivityID,
			ScannedAt:  scannedAt,
			Status:     "checked",
			ScannedBy:  scannedBy,
		})
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
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check existing check-in")
		return
	}

	user, err := h.DB.GetCheckInLog(id)
	if err != nil {
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch existing check-in")
		return
	}
	if user.Status == "checked" {
		utils.RespondWithError(w, http.StatusConflict, "Cannot Check in twice")
		return
	}

	user.Status = "checked"
	user.ScannedBy = fbuser.Email
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

	checkIn, err := h.DB.GetCheckInLog(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "check in id not found")
		return
	}

	user, err := h.DB.GetUser(checkIn.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "user id not found")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "checkedIn in id not found")
		return
	}

	if checkIn.Status == "checked" {
		checkIn.Status = "unchecked"
		err := h.DB.UpdateCheckInLog(checkIn)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
			return
		}

		checkInReponse := models.CheckInRespose{
			ID:         checkIn.ID,
			UserID:     checkIn.UserID,
			ActivityID: checkIn.ActivityID,
			Status:     checkIn.Status,
			FullName:   user.FullName,
			ScannedAt:  checkIn.ScannedAt,
			ScannedBy:  checkIn.ScannedBy,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(checkInReponse)
		return
	}

	checkIn.Status = "checked"
	err = h.DB.UpdateCheckInLog(checkIn)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	checkInReponse := models.CheckInRespose{
		ID:         checkIn.ID,
		UserID:     checkIn.UserID,
		ActivityID: checkIn.ActivityID,
		Status:     checkIn.Status,
		FullName:   user.FullName,
		ScannedAt:  checkIn.ScannedAt,
		ScannedBy:  checkIn.ScannedBy,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkInReponse)
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Can't get check-in logs")
		return
	}

	var responses []models.CheckInRespose

	for _, logItem := range checkInLogs {
		user, err := h.DB.GetUser(logItem.UserID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Can't get user details")
			return
		}
		activity, err := h.DB.GetActivity(logItem.ActivityID)
		if err != nil {
			resp := models.CheckInRespose{
				ID:           logItem.ID,
				FullName:     user.FullName,
				UserID:       logItem.UserID,
				ActivityName: "Cant Find",
				ActivityID:   logItem.ActivityID,
				ScannedAt:    logItem.ScannedAt,
				ScannedBy:    logItem.ScannedBy,
				Status:       logItem.Status,
			}
			responses = append(responses, resp)
			continue
		}

		resp := models.CheckInRespose{
			ID:           logItem.ID,
			FullName:     user.FullName,
			UserID:       logItem.UserID,
			ActivityName: activity.Name,
			ActivityID:   logItem.ActivityID,
			ScannedAt:    logItem.ScannedAt,
			ScannedBy:    logItem.ScannedBy,
			Status:       logItem.Status,
		}
		responses = append(responses, resp)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

/*
Export CheckIn , retrives all check Ins

Returns:
- 200 OK with updated check-in JSON on success
- 500 Internal Server Error on DB failure
*/
func (h *Handler) ExportCheckIn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["event_id"]
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	checkInLogs, err := h.DB.GetAllCheckInOfEvents(id)
	if err != nil {
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Can't get check-in logs")
		return
	}

	f := excelize.NewFile()
	sheet := "CheckIns"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Full Name", "Activity ", "Scanned At", "Scanned By", "Status"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, header)
	}

	for i, logItem := range checkInLogs {
		user, err := h.DB.GetUser(logItem.UserID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Can't get user details")
			return
		}
		activity, err := h.DB.GetActivity(logItem.ActivityID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Can't get user details")
			return
		}

		rowNum := i + 2 // excel rows start at 1 and row 1 is the header so manually +2

		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), user.AutoId)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), user.FullName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), activity.Name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), logItem.ScannedAt.Format(time.RFC3339))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), logItem.ScannedBy)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), logItem.Status)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="checkins.xlsx"`)
	w.WriteHeader(http.StatusOK)

	err = f.Write(w)
	if err != nil {
		log.Printf("Error writing Excel file: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to write Excel file")
	}

}

/*
GetCheckInById , retrives all check Ins

Returns:
- 200 OK with updated check-in JSON on success
- 500 Internal Server Error on DB failure
*/

func (h *Handler) GetCheckInByEventId(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing GetCheckInByEventId handler")
	vars := mux.Vars(r)
	eventIdStr := vars["event_id"]
	if eventIdStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	event_id, err := uuid.Parse(eventIdStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	checkInLogs, err := h.DB.GetAllCheckInOfEvents(event_id)
	if err != nil {
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Can't get check-in logs")
		return
	}

	var responses []models.CheckInRespose

	for _, logItem := range checkInLogs {
		user, err := h.DB.GetUser(logItem.UserID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Can't get user details")
			return
		}

		activity, err := h.DB.GetActivity(logItem.ActivityID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Can't get activity name"+err.Error())
			return
		}

		resp := models.CheckInRespose{
			ID:           logItem.ID,
			FullName:     user.FullName,
			UserID:       logItem.UserID,
			ActivityName: activity.Name,
			ActivityID:   logItem.ActivityID,
			ScannedAt:    logItem.ScannedAt,
			ScannedBy:    logItem.ScannedBy,
			Status:       logItem.Status,
		}

		responses = append(responses, resp)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

/*
GetCheckInById , retrives all check Ins

Returns:
- 200 OK with updated check-in JSON on success
- 500 Internal Server Error on DB failure
*/

func (h *Handler) GetCheckInByActivityId(w http.ResponseWriter, r *http.Request) {

	fbUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: no user in context")
		return
	}
	vars := mux.Vars(r)

	activityIdStr := vars["activity_id"]
	if activityIdStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	activityId, err := uuid.Parse(activityIdStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
  
  eventId , err := h.DB.GetEventIdByActivity(activityId)
  if err != nil {
    log.Println("here1")
    log.Println(err.Error())
    utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get event ID by activity")
    return
  }
  
	access, err := h.DB.CanSeeScanned(fbUser.UID, eventId.String())

  if err != nil {
    log.Println("here2")
    log.Println(err.Error())
    utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check event access")
    return
  }
  
  if !access {
    utils.RespondWithError(w, http.StatusUnauthorized, "Access denied")
    return
  }

	checkInLogs, err := h.DB.GetAllCheckInOfActivity(activityId)
	if err != nil {
    log.Println("here3")
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Can't get check-in logs")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkInLogs)
}

/*
GetCheckInById , retrives all check Ins

Returns:
- 200 OK with updated check-in JSON on success
- 500 Internal Server Error on DB failure
*/

func (h *Handler) GetCheckInByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIdStr := vars["attendee_id"]
	if activityIdStr == "" {
		log.Print("Missing ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	activityId, err := uuid.Parse(activityIdStr)
	if err != nil {
		log.Print("Invalid ID format: ", err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	checkInLogs, err := h.DB.GetAllCheckInOfUser(activityId)
	if err != nil {
		log.Print(err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Can't get check-in logs")
		return
	}
	log.Print(checkInLogs)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkInLogs)
}
