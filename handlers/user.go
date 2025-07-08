package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
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

	if u.Email == "" || u.FullName == "" || u.Phone == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}


	err := h.DB.CreateUser(&u)

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
	json.NewEncoder(w).Encode(u)
}
