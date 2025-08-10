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

  log.Print("here")
	connStr := os.Getenv("POSTGRESS_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	Router := mux.NewRouter()
	Router.Use(fbAuth.AuthMiddleware)

	db, err := postgres.ConnectPostgres(connStr)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("database count not be connected")
	}

  log.Print("here")
	handler := handlers.New(db, fbAuth)

	Router.HandleFunc("/user", handler.CreateUser).Methods(constants.Post)
	Router.HandleFunc("/users/{event_id}", handler.GetUsersByEvent).Methods(constants.Get)
	Router.HandleFunc("/importusers/{event_id}", handler.ImportUser).Methods(constants.Post)

	Router.HandleFunc("/event", handler.CreateEvent).Methods(constants.Post)
	Router.HandleFunc("/event", handler.GetEvent).Methods(constants.Get)
	Router.HandleFunc("/eventinfo", handler.GetEventInfo).Methods(constants.Get)
	Router.HandleFunc("/addeventwithcode/{code}", handler.AddEventWithEventCode).Methods(constants.Post)
	Router.HandleFunc("/giveRoleToStaffs", handler.GiveRoleToStaff).Methods(constants.Post)
	Router.HandleFunc("/modifyRoleToStaffs", handler.ModifyRoleToStaff).Methods(constants.Post)
	Router.HandleFunc("/getstaffs/{event_id}", handler.GetStaffsByEvent).Methods(constants.Get)

	Router.HandleFunc("/activity", handler.CreateActivity).Methods(constants.Post)

	Router.HandleFunc("/checkins", handler.GetCheckIn).Methods(constants.Get)
	Router.HandleFunc("/checkins/{event_id}", handler.GetCheckInByEventId).Methods(constants.Get)
	Router.HandleFunc("/activitycheckins/{activity_id}", handler.GetCheckInByActivityId).Methods(constants.Get)
	Router.HandleFunc("/attendeecheckins/{attendee_id}", handler.GetCheckInByUserId).Methods(constants.Get)
	Router.HandleFunc("/checkins", handler.CreateCheckIn).Methods(constants.Post)
  Router.HandleFunc("/checkins/{id}", handler.ModifyCheckIn).Methods(constants.Put)
	Router.HandleFunc("/exportcheckins/{event_id}", handler.ExportCheckIn).Methods(constants.Get)

	log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, WithCORS(Router))

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
