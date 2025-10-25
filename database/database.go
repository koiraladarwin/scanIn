package db

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

type Database interface {
	CreateUserCategory(*models.UsersCategoryRequest) (*models.UsersCategory, error)
	GetUserCategories(firebaseId string) ([]models.UsersCategory, error)
	CreateUser(*models.UserRequest) (*models.User, error)
	GetUsers(firebaseid string) ([]models.User, error)
	GetUser(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.UserModifyRequest) error
	GetUsersByEvent(eventID uuid.UUID) ([]models.User, error)
	DeleteUser(id uuid.UUID) error

	CreateEventCategory(*models.EventCategoryRequest) (models.EventCategory, error)
	GetEventsWithDetails(firebaseId string) ([]models.EventWithDetails, error)
	GetEventsWithSessions(firebaseId string) ([]models.EventwithSessions, error)
	GetEventCategories(firebase_id string) ([]models.EventCategory, error)
	CreateEvent(*models.EventCreateRequest) (models.Event, error)
	UpdateEvent(*models.EventModifyRequest) error
	DeleteEvent(id uuid.UUID) error
	EventExists(eventID uuid.UUID) (bool, error)
	GetAllEvents() ([]models.Event, error)
	GetEventsByFirebaseUser(firebaseId string) ([]models.Event, error)
	GetEventByFirebaseUser(firebaseId string, eventId uuid.UUID) (*models.Event, error)
	GetEventByAdminId(id string) (*models.Event, error)
	GetEventByStaffId(id string) (*models.Event, error)

	CreateActivity(*models.ActivityCreateRequest) (*models.ActivityWithScannedUser, error)
	GetActivitiesDetails(firebaseId string) ([]models.ActivityDetails, error)
	GetActivity(id uuid.UUID) (*models.ActivityWithScannedUser, error)
	UpdateActivity(*models.ActivityWithScannedUser) error
	DeleteActivity(id uuid.UUID) error
	GetActivitiesByEvent(firebaseId string, eventID uuid.UUID) ([]models.ActivityWithScannedUser, error)
	GetEventIdByActivity(activityId uuid.UUID) (uuid.UUID, error)

	CreateCheckInLog(*models.CheckInLog) (*models.CheckInLog, error)
	GetCheckInLog(id uuid.UUID) (*models.CheckInLog, error)
	DeleteCheckInLog(id uuid.UUID) error
	CheckInExists(userID uuid.UUID, activityID uuid.UUID) (uuid.UUID, error)
	GetAllCheckInLog() ([]models.CheckInLog, error)
	GetAllCheckInOfEvents(eventID uuid.UUID) ([]models.CheckInLog, error)
	GetAllCheckInOfActivity(activityID uuid.UUID) ([]models.CheckInRespose, error)
	GetAllCheckInOfUser(userID uuid.UUID) ([]models.CheckInRespose, error)

	CreateAttendee(a models.AttendeeRequest) (models.Attendee, error)
	GetAttendeeFromUserIdandTicketId(userID, ticketID uuid.UUID) (models.Attendee, error)
	CreateAttendeeActivityEnroll(a models.AttendeeActivityRequest) (models.AttendeeActivity, error)
	GetEventFromAttendee(attendeeID uuid.UUID) (models.Event, error)
	GetEnrolledAttendee(firebaseID string) ([]models.EnrolledAttendee, error)

	GetTicketsForAllEvents(firebaseId string) ([]models.EventsTicket, error)
	GetInviteeForAllEvents(firebaseId string) ([]models.EventsInvitee, error)
	GetEventsWithSessionsAndTickets(firebaseId string) ([]models.EventwithSessionsAndTickets, error)
	CreateTicket(a models.TicketRequest) (models.Ticket, error)
	GetTicketById(firebaseId string, id uuid.UUID) (models.Ticket, error)
	CreateTicketCategory(a models.TicketCategory) (models.TicketCategory, error)
	GetTicketCategories(firebaseId string, ticket_type string) ([]models.TicketCategory, error)
	GetTicketCategory(firebaseId string, id uuid.UUID) (models.TicketCategory, error)

	CreateStaffCategory(staffCategoryRequest *models.StaffCategory) (*models.StaffCategory, error)
	GetStaffCategories(firebaseId string) ([]models.StaffCategory, error)
	CreateStaff(staffRequest *models.Staff) (*models.Staff, error)
	GetStaffs(firebaseId string) ([]models.Staff, error)
	CreateStaffEventEnroll(staffEnrollRequest *models.StaffEnroll) (*models.StaffEnroll, error)
	GetStaffEventEnroll(firebaseId string,  staffId uuid.UUID,eventId uuid.UUID,) (models.StaffEnroll, error)
	CreateStaffActivityAssign(staffActivityRequest *models.StaffActivities) (*models.StaffActivities, error)
	GetEnrolledStaff(firebaseID string) ([]models.EnrolledStaff, error)
	Close() error
}
