package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	handler := handlers.New(db)

	api.HandleFunc("/user", handler.CreateUser).Methods(post)
	
  api.HandleFunc("/attendees", handler.RegisterAttendee).Methods(post)
  api.HandleFunc("/attendees/{event_id}", handler.GetAttendeesByEvent).Methods(get)
	
  api.HandleFunc("/checkins", handler.CreateCheckIn).Methods(post)
	api.HandleFunc("/checkins/{id}", handler.ModifyCheckIn).Methods(put)
	api.HandleFunc("/checkins/{id}", handler.ModifyCheckIn).Methods(put)

	log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

