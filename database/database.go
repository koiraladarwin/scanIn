package db

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

type Database interface {
	CreateUser(*models.UserRequest) (*models.User, error)
	GetUser(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.UserModifyRequest)  (error)
	GetUsersByEvent(eventID uuid.UUID) ([]models.User, error)

	CreateEvent(*models.EventCreateRequest) error
	UpdateEvent(*models.EventModifyRequest) error
	DeleteEvent(id uuid.UUID) error
	EventExists(eventID uuid.UUID) (bool, error)
	GetAllEvents() ([]models.Event, error)
	GetEventsByFirebaseUser(firebaseId string) ([]models.Event, error)
	GetEventByFirebaseUser(firebaseId string, eventId uuid.UUID) (*models.Event, error)
	GetEventByAdminId(id string) (*models.Event, error)
	GetEventByStaffId(id string) (*models.Event, error)
  GetStaffByEvent(eventId string) ([]models.Staff, error)

	CreateActivity(*models.ActivityCreateRequest) error
	GetActivity(id uuid.UUID) (*models.Activity, error)
	UpdateActivity(*models.Activity) error
	DeleteActivity(id uuid.UUID) error
	GetActivitiesByEvent(firebaseId string, eventID uuid.UUID) ([]models.Activity, error)
	GetEventIdByActivity(activityId uuid.UUID) (uuid.UUID, error)

	CreateCheckInLog(*models.CheckInLog) error
	GetCheckInLog(id uuid.UUID) (*models.CheckInLog, error)
	UpdateCheckInLog(*models.CheckInLog) error
	DeleteCheckInLog(id uuid.UUID) error
	CheckInExists(userID uuid.UUID, activityID uuid.UUID) (uuid.UUID, error)
	GetAllCheckInLog() ([]models.CheckInLog, error)
	GetAllCheckInOfEvents(eventID uuid.UUID) ([]models.CheckInLog, error)
	GetAllCheckInOfActivity(activityID uuid.UUID) ([]models.CheckInRespose, error)
	GetAllCheckInOfUser(userID uuid.UUID) ([]models.CheckInRespose, error)

	IsCreator(fbId string, eventId string) (bool, error)
	CanSeeScanned(fbId string, eventId string) (bool, error)
	CanCreateActivity(fbId string, eventId string) (bool, error)
	CanCreateAttendee(fbId string, eventId string) (bool, error)
	CanSeeAttendee(fbId string, eventId string) (bool, error)
	CanSeeEventInfo(fbId, eventId string) (bool, error)
	AddStaffToEvent(fbId, eventId string) error
	AddAdminToEvent(fbId, eventId string) error
	AddEventRole(role models.RoleRequest) error
	ModifyEventRole(role models.EditRoleRequest) error

	Close() error
}
