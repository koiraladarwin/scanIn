package db

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

type Database interface {
	CreateUser(*models.UserRequest) (*models.User, error)
	GetUser(id uuid.UUID) (*models.User, error)
	GetUsersByEvent(eventID uuid.UUID) ([]models.User, error)

	CreateEvent(*models.EventRequest) error
	GetEvent(id uuid.UUID) (*models.Event, error)
	UpdateEvent(*models.Event) error
	DeleteEvent(id uuid.UUID) error
	EventExists(eventID uuid.UUID) (bool, error)
	GetAllEvents() ([]models.Event, error)

	CreateActivity(*models.ActivityRequest) error
	GetActivity(id uuid.UUID) (*models.Activity, error)
	UpdateActivity(*models.Activity) error
	DeleteActivity(id uuid.UUID) error
	GetActivitiesByEvent(eventID uuid.UUID) ([]models.Activity, error)

	CreateAttendee(*models.Attendee) (*models.Attendee, error)
	GetAttendee(id uuid.UUID) (*models.Attendee, error)
	UpdateAttendee(*models.Attendee) error
	DeleteAttendee(id uuid.UUID) error
	GetAttendeesByEvent(eventID uuid.UUID) ([]models.Attendee, error)
	GetNumberOfUsersByEvent(eventID uuid.UUID) (int, error)

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

	Close() error
}
