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

	CreateActivity(*models.ActivityCreateRequest) (*models.Activity, error)
	GetActivity(id uuid.UUID) (*models.Activity, error)
	UpdateActivity(*models.Activity) error
	DeleteActivity(id uuid.UUID) error
	GetActivitiesByEvent(firebaseId string, eventID uuid.UUID) ([]models.Activity, error)
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

	CreateTicket(a models.TicketRequest) (models.Ticket, error)
	CreateTicketCategory(a models.TicketCategory) (models.TicketCategory, error)
	GetTicketCategories(firebaseId string, ticket_type string) ([]models.TicketCategory, error)
	GetTicketCategory(firebaseId string, id uuid.UUID) (models.TicketCategory, error)

	CreateStaffCategory(staffCategoryRequest *models.StaffCategory) (*models.StaffCategory, error)
	GetStaffCategories(firebaseId string) ([]models.StaffCategory, error)
	Close() error
	CreateStaff(staffRequest *models.Staff) (*models.Staff, error)
	GetStaffs(firebaseId string) ([]models.Staff, error)

}
