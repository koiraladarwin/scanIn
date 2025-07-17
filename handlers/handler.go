package handlers

import (
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
)

type Handler struct {
	DB     db.Database
	FbAuth *firebaseauth.FirebaseAuth
}

func New(db db.Database,fbAuth *firebaseauth.FirebaseAuth) *Handler {
	return &Handler{DB: db,FbAuth: fbAuth}
}
