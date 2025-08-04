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
	"github.com/koiraladarwin/scanin/handlers/middleware"
)

const (
	AdminAccessLevel = 2
	StaffAccessLevel = 1
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

	AdminRouter := mux.NewRouter()
	AdminRouter.Use(fbAuth.AuthMiddleware)

	db, err := postgres.ConnectPostgres(connStr)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("database count not be connected")
	}

	handler := handlers.New(db, fbAuth)

	AdminRouter.HandleFunc("/user", middleware.RequireAccessLevel(AdminAccessLevel, handler.CreateUser)).Methods(constants.Post)
	AdminRouter.HandleFunc("/users/{event_id}", middleware.RequireAccessLevel(StaffAccessLevel, handler.GetUsersByEvent)).Methods(constants.Get)
	AdminRouter.HandleFunc("/importusers/{event_id}", middleware.RequireAccessLevel(AdminAccessLevel, handler.ImportUser)).Methods(constants.Post)

	AdminRouter.HandleFunc("/event", middleware.RequireAccessLevel(AdminAccessLevel, handler.CreateEvent)).Methods(constants.Post)
	AdminRouter.HandleFunc("/event", middleware.RequireAccessLevel(StaffAccessLevel, handler.GetEvent)).Methods(constants.Get)
	AdminRouter.HandleFunc("/eventinfo", middleware.RequireAccessLevel(StaffAccessLevel, handler.GetEventInfo)).Methods(constants.Get)

	AdminRouter.HandleFunc("/activity", middleware.RequireAccessLevel(AdminAccessLevel, handler.CreateActivity)).Methods(constants.Post)

	AdminRouter.HandleFunc("/checkins", middleware.RequireAccessLevel(StaffAccessLevel, handler.GetCheckIn)).Methods(constants.Get)
	AdminRouter.HandleFunc("/checkins/{event_id}", middleware.RequireAccessLevel(StaffAccessLevel, handler.GetCheckInByEventId)).Methods(constants.Get)
	AdminRouter.HandleFunc("/activitycheckins/{activity_id}", middleware.RequireAccessLevel(StaffAccessLevel, handler.GetCheckInByActivityId)).Methods(constants.Get)
	AdminRouter.HandleFunc("/checkins", middleware.RequireAccessLevel(StaffAccessLevel, handler.CreateCheckIn)).Methods(constants.Post)
	AdminRouter.HandleFunc("/checkins/{id}", middleware.RequireAccessLevel(AdminAccessLevel, handler.ModifyCheckIn)).Methods(constants.Put)
	AdminRouter.HandleFunc("/exportcheckins/{event_id}", middleware.RequireAccessLevel(AdminAccessLevel, handler.ExportCheckIn)).Methods(constants.Get)

	log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, WithCORS(AdminRouter))

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
