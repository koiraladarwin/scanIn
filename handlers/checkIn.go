package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
	"github.com/xuri/excelize/v2"
)

// todo: make sure the firebaseuser has right to create the checkin
func (h *Handler) CreateCheckIn(w http.ResponseWriter, r *http.Request) {
	firebaseUser, ok := firebaseauth.FbUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: no user in context")
		return
	}

	var checkInReq models.CheckInLogRequest
	err := json.NewDecoder(r.Body).Decode(&checkInReq)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if checkInReq.AttendeeId == uuid.Nil || checkInReq.ActivityID == uuid.Nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	checkIn := &models.CheckInLog{
		AttendeeId: checkInReq.AttendeeId,
		ActivityID: checkInReq.ActivityID,
		ScannedBy:  firebaseUser.UID,
		ScannedAt:  time.Now(),
	}

	checkIn, err = h.DB.CreateCheckInLog(checkIn)

	if err != nil {
		log.Println("Error creating check-in log:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create check-in log")
		return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

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
		user, err := h.DB.GetUser(logItem.AttendeeId)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Can't get user details")
			return
		}
		activity, err := h.DB.GetActivity(logItem.ActivityID)
		if err != nil {
			resp := models.CheckInRespose{
				ID:           logItem.ID,
				FullName:     user.FullName,
				AttendeeId:   logItem.AttendeeId,
				ActivityName: "Cant Find",
				ActivityID:   logItem.ActivityID,
				ScannedAt:    logItem.ScannedAt,
				ScannedBy:    logItem.ScannedBy,
			}
			responses = append(responses, resp)
			continue
		}

		resp := models.CheckInRespose{
			ID:           logItem.ID,
			FullName:     user.FullName,
			AttendeeId:   logItem.AttendeeId,
			ActivityName: activity.Name,
			ActivityID:   logItem.ActivityID,
			ScannedAt:    logItem.ScannedAt,
			ScannedBy:    logItem.ScannedBy,
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
		user, err := h.DB.GetUser(logItem.AttendeeId)
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
		user, err := h.DB.GetUser(logItem.AttendeeId)
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
			AttendeeId:   logItem.AttendeeId,
			ActivityName: activity.Name,
			ActivityID:   logItem.ActivityID,
			ScannedAt:    logItem.ScannedAt,
			ScannedBy:    logItem.ScannedBy,
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

	_, ok := firebaseauth.FbUserFromContext(r.Context())
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
- 400 Bad Request
*/

func (h *Handler) GetCheckInByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkInLogs := []models.CheckInRespose{}
	activityIdStr := vars["attendee_id"]
	if activityIdStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	activityId, err := uuid.Parse(activityIdStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	checkInLogs, err = h.DB.GetAllCheckInOfUser(activityId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Can't get check-in logs")
		return
	}
	log.Print(checkInLogs)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkInLogs)
}
