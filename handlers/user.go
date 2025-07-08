package handlers
import (
	"encoding/json"
	"net/http"

	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

/*
CreateUser accepts JSON:
{
  "full_name": "string",
  "email": "string",
  "phone": "string",
  "role": "participant|staff|member"
}

Returns:
- 201 Created with created user JSON on success
- 400 Bad Request for invalid input
- 500 Internal Server Error on DB failure
*/
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := h.DB.CreateUser(&u); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

