package handlers

import (
	"github.com/koiraladarwin/scanin/database"
)

type Handler struct {
	DB db.Database
}

func New(db db.Database) *Handler {
	return &Handler{DB: db}
}
