package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/koiraladarwin/scanin/constants"
	"github.com/koiraladarwin/scanin/database/postgres"
	"github.com/koiraladarwin/scanin/features/firebaseauth"
	"github.com/koiraladarwin/scanin/handlers"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*clientLimiter)
	mu      sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	cl, exists := clients[ip]
	if !exists {
		cl = &clientLimiter{
			limiter:  rate.NewLimiter(1, 10),
			lastSeen: time.Now(),
		}
		clients[ip] = cl
	} else {
		cl.lastSeen = time.Now()
	}
	return cl.limiter
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func timeoutMiddleware(timeout time.Duration) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func withCORS(next http.Handler) http.Handler {
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

	Router := mux.NewRouter()
	Router.Use(rateLimitMiddleware)
	Router.Use(fbAuth.AuthMiddleware)
	Router.Use(timeoutMiddleware(5 * time.Second))

	db, err := postgres.ConnectPostgres(connStr)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("database count not be connected")
	}

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

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, withCORS(Router))

	if err != nil {
		log.Fatal(err)
	}
}
