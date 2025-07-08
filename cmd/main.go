package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/constants"
	"github.com/koiraladarwin/scanin/database/postgres"
	"github.com/koiraladarwin/scanin/handlers"
)


func main() {
	port := "4000"
	connStr := "postgres://postgres:mysecretpassword@localhost:5432/scanin?sslmode=disable"

	db, err := postgres.ConnectPostgres(connStr)
	if err != nil {
		log.Fatal("could not connect to database:", err)
	}

	api := mux.NewRouter()

	handler := handlers.New(db)

	api.HandleFunc("/user", handler.CreateUser).Methods(constants.Post)
  api.HandleFunc("/users/{event_id}", handler.GetUsersByEvent).Methods(constants.Get)
	
  api.HandleFunc("/event", handler.CreateEvent).Methods(constants.Post)
  
  api.HandleFunc("/activity", handler.CreateActivity).Methods(constants.Post)
  
  api.HandleFunc("/attendees", handler.RegisterAttendee).Methods(constants.Post)
  
  api.HandleFunc("/checkins", handler.CreateCheckIn).Methods(constants.Post)
	api.HandleFunc("/checkins/{id}", handler.ModifyCheckIn).Methods(constants.Put)

	
  log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, api)
	if err != nil {
		log.Fatal(err)
	}
}

