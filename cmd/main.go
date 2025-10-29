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

	ctx := context.Background()
	if os.Getenv("RAILWAY_ENVIRONMENT_ID") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, using environment variables instead")
		}
	}

	port := os.Getenv("PORT") 
	if port == "" {
		port = "4000" 
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
	Router.Use(fbAuth.AuthMiddleware)

	db, err := postgres.ConnectPostgres(connStr)
	if err != nil {
		log.Print(err.Error())
		log.Fatal("database count not be connected")
	}

	handler := handlers.New(db, fbAuth)
  
  Router.HandleFunc("/attendeecategory", handler.CreateUserCategory).Methods(constants.Post)
  Router.HandleFunc("/attendeecategories", handler.GetUserCategories).Methods(constants.Get)
	Router.HandleFunc("/attendee", handler.CreateUser).Methods(constants.Post)
	Router.HandleFunc("/attendee", handler.GetUsers).Methods(constants.Get)
	Router.HandleFunc("/modifyattendee", handler.UpdateUser).Methods(constants.Put)
	Router.HandleFunc("/deleteattendee", handler.DeleteUser).Methods(constants.Delete)
	Router.HandleFunc("/attendee/{event_id}", handler.GetUsersByEvent).Methods(constants.Get)
  
  Router.HandleFunc("/enrollattendeeevent", handler.CreateAttendee).Methods(constants.Post)
  Router.HandleFunc("/enrollattendeeactivity", handler.CreateAttendeeActivityEnroll).Methods(constants.Post)
  Router.HandleFunc("/enrollattendee", handler.EnrollAttendee).Methods(constants.Post)
  Router.HandleFunc("/enrollattendee", handler.GetEnrolledAttendee).Methods(constants.Get)
  Router.HandleFunc("/ticketattendee", handler.GetTicketAttendee).Methods(constants.Get) 
  Router.HandleFunc("/ticketstatus", handler.ChangeTicketAttendeeStatus).Methods(constants.Put) 

  Router.HandleFunc("/ticket", handler.CreateTicket).Methods(constants.Post) 
  Router.HandleFunc("/eventtickets", handler.GetTicketsForAllEvents).Methods(constants.Get) 
  Router.HandleFunc("/ticketcategory", handler.CreateTicketCategory).Methods(constants.Post) 
  Router.HandleFunc("/ticketcategories", handler.GetTicketCategories).Methods(constants.Get) 
  
  Router.HandleFunc("/eventinvitees", handler.GetInviteeForAllEvents).Methods(constants.Get) 
  Router.HandleFunc("/invitee", handler.CreateInvitee).Methods(constants.Post)
  Router.HandleFunc("/inviteecategory", handler.CreateInviteeCategory).Methods(constants.Post)
  Router.HandleFunc("/inviteecategories", handler.GetInviteeCategories).Methods(constants.Get)

  Router.HandleFunc("/eventcategory", handler.CreateEventCategory).Methods(constants.Post)  
  Router.HandleFunc("/eventcategories", handler.GetEventCategories).Methods(constants.Get)
	Router.HandleFunc("/event", handler.CreateEvent).Methods(constants.Post)
	Router.HandleFunc("/event", handler.GetEvent).Methods(constants.Get)
	Router.HandleFunc("/eventswithdetails", handler.GetEventsWithDetails).Methods(constants.Get)
	Router.HandleFunc("/eventswithsessions", handler.GetEventsWithSesion).Methods(constants.Get)
	Router.HandleFunc("/eventswithsessionsandtickets", handler.GetEventsWithSesionsAndTickets).Methods(constants.Get)
	Router.HandleFunc("/modifyevent", handler.ModifyEvent).Methods(constants.Put)
	Router.HandleFunc("/eventinfo", handler.GetEventInfo).Methods(constants.Get)
	Router.HandleFunc("/addeventwithcode/{code}", handler.AddEventWithEventCode).Methods(constants.Post)

	Router.HandleFunc("/session", handler.CreateActivity).Methods(constants.Post)
	Router.HandleFunc("/sessionwithdetails", handler.GetAcivityWithDetails).Methods(constants.Get)
	Router.HandleFunc("/modifyactivity", handler.UpdateActivity).Methods(constants.Put)

	Router.HandleFunc("/checkins", handler.GetCheckIn).Methods(constants.Get)
	Router.HandleFunc("/checkins/{event_id}", handler.GetCheckInByEventId).Methods(constants.Get)
	Router.HandleFunc("/activitycheckins/{activity_id}", handler.GetCheckInByActivityId).Methods(constants.Get)
	Router.HandleFunc("/attendeecheckins/{attendee_id}", handler.GetCheckInByUserId).Methods(constants.Get)
	Router.HandleFunc("/checkins", handler.CreateCheckIn).Methods(constants.Post)
	Router.HandleFunc("/exportcheckins/{event_id}", handler.ExportCheckIn).Methods(constants.Get)
  
  Router.HandleFunc("/staffcategory", handler.CreateStaffCategory).Methods(constants.Post)
  Router.HandleFunc("/staffcategories", handler.GetStaffCategories).Methods(constants.Get)
  Router.HandleFunc("/staff", handler.CreateStaff).Methods(constants.Post)
  Router.HandleFunc("/staff", handler.GetStaffs).Methods(constants.Get)
  Router.HandleFunc("/enrollstaff", handler.CreateStaffEnrollment).Methods(constants.Post)
  Router.HandleFunc("/enrollstaff", handler.GetEnrolledStaff).Methods(constants.Get)

	log.Printf("Server running on port %s", port)
	err = http.ListenAndServe(":"+port, withCORS(Router))

	if err != nil {
		log.Fatal(err)
	}
}
