package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/koiraladarwin/scanin/constants"
	"github.com/koiraladarwin/scanin/database/postgres"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/handlers"
)

func main() {
	port := "4000"
	ctx := context.Background()
	if os.Getenv("RAILWAY_ENVIRONMENT_ID") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, using environment variables instead")
		}
	}

	fbAuth, err := firebaseauth.NewFirebaseAuth(ctx)
	if err != nil {
		log.Fatal("firebase Auth could not be instatitated")
	}
  

	connStr := os.Getenv("POSTGRESS_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	api := mux.NewRouter()

	_ = api.NewRoute().Subrouter()
	db, err := postgres.ConnectPostgres(connStr)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("database count not be connected")
	}



	handler := handlers.New(db, fbAuth)

  api.HandleFunc("/user", handler.CreateUser).Methods(constants.Post)
	api.HandleFunc("/user", handler.GetUser).Methods(constants.Get)
	api.HandleFunc("/users/{event_id}", handler.GetUsersByEvent).Methods(constants.Get)

	api.HandleFunc("/event", handler.CreateEvent).Methods(constants.Post)
	api.HandleFunc("/event", handler.GetEvent).Methods(constants.Get)
	api.HandleFunc("/eventinfo", handler.GetEventInfo).Methods(constants.Get)

	api.HandleFunc("/activity", handler.CreateActivity).Methods(constants.Post)

	api.HandleFunc("/attendees", handler.RegisterAttendee).Methods(constants.Post)

	api.HandleFunc("/checkins", handler.GetCheckIn).Methods(constants.Get)
	api.HandleFunc("/checkins/{id}", handler.GetCheckInById).Methods(constants.Get)
	api.HandleFunc("/checkins", handler.CreateCheckIn).Methods(constants.Post)
	api.HandleFunc("/checkins/{id}", handler.ModifyCheckIn).Methods(constants.Put)
	api.HandleFunc("/exportcheckins/{event_id}", handler.ExportCheckIn).Methods(constants.Get)

	log.Printf("Server running on port %s", port)

	err = http.ListenAndServe(":"+port, WithCORS(fbAuth.AuthMiddleware(api)))

	if err != nil {
		log.Fatal(err)
	}
}

func WithCORS(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    if r.Method == http.MethodOptions {
      w.WriteHeader(http.StatusNoContent)
      return
    }
    next.ServeHTTP(w, r)
  })
}

