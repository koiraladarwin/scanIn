package main

import (
	"github.com/gorilla/mux"
	"github.com/koiraladarwin/scanin/constants"
	"github.com/koiraladarwin/scanin/database/postgres"
	"github.com/koiraladarwin/scanin/handlers"
	"log"
	"net/http"
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
	api.HandleFunc("/user", handler.GetUser).Methods(constants.Get)
	api.HandleFunc("/users/{event_id}", handler.GetUsersByEvent).Methods(constants.Get)

	api.HandleFunc("/event", handler.CreateEvent).Methods(constants.Post)
	api.HandleFunc("/event", handler.GetEvent).Methods(constants.Get)
	api.HandleFunc("/eventinfo", handler.GetEventInfo).Methods(constants.Get)

	api.HandleFunc("/activity", handler.CreateActivity).Methods(constants.Post)

	api.HandleFunc("/attendees", handler.RegisterAttendee).Methods(constants.Post)

	api.HandleFunc("/checkins", handler.GetCheckIn).Methods(constants.Get)
	api.HandleFunc("/checkins", handler.CreateCheckIn).Methods(constants.Post)
	api.HandleFunc("/checkins/{id}", handler.ModifyCheckIn).Methods(constants.Put)

	log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, WithCORS(api))
	if err != nil {
		log.Fatal(err)
	}
}

// temporary only during development
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
